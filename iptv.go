package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/hauke96/sigolo"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	FullIPTV   = "full_iptv.m3u"
	ParsedIPTV = "iptv.m3u"
)

func processM3u8(conf *Config) error {
	sigolo.Info("processM3u8")
	// get last modified time
	file, err := os.Stat(FullIPTV)

	var (
		firstRun bool
		timeDiff float64
	)

	if os.IsNotExist(err) {
		firstRun = true
		_, err = os.Create(FullIPTV)
		if err != nil {
			return err
		}
	}

	if firstRun {
		timeDiff = 99999999
	} else {
		timeDiff = time.Now().Sub(file.ModTime()).Minutes()
	}

	// if file is older than 2 hours
	if timeDiff > float64(conf.CacheTime) {
		sigolo.Info("getting new iptv file")
		err := getNewFullFile(conf)
		if err != nil {
			return err
		}

		err = parseFullIPTV(conf)
		if err != nil {
			return err
		}
	} else {
		sigolo.Info("skip getting new file, valid for %f minute(s)", float64(conf.CacheTime)-timeDiff)
	}

	return nil
}

func getNewFullFile(conf *Config) error {
	sigolo.Info("getNewFullFile")

	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(FullIPTV + ".tmp")
	if err != nil {
		return err
	}

	resp, err := http.Get(conf.Target)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Print("\n")

	// Close the file without defer so it can happen before Rename()
	out.Close()

	if err = os.Rename(FullIPTV+".tmp", FullIPTV); err != nil {
		return err
	}
	return nil
}

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

func parseFullIPTV(conf *Config) error {
	sigolo.Info("parseFullIPTV")

	channels, err := Parse(FullIPTV)
	if err != nil {
		return err
	}

	var copyGroup Channels
	for _, v := range channels {
		for _, g := range conf.Include {
			if v.GroupName == g {
				copyGroup = append(copyGroup, v)
			}
		}
	}

	err = writeLines(conf, copyGroup, ParsedIPTV)
	if err != nil {
		return err
	}

	return nil
}
