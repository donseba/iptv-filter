package main

import (
	"github.com/hauke96/sigolo"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	FullIPTV   = "full_iptv.m3u"
	ParserIPTV = "iptv.m3u"
)

func processM3u8(conf *Config) error {
	sigolo.Info("processM3u8")
	// get last modified time
	file, err := os.Stat(FullIPTV)

	if err != nil {
		return err

	}

	timeDiff := time.Now().Sub(file.ModTime()).Minutes()

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

	resp, err := http.Get(conf.Target)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	out, err := os.Create(FullIPTV)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
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

	err = writeLines(conf, copyGroup, ParserIPTV)
	if err != nil {
		return err
	}

	return nil
}
