package main

import (
	"bufio"
	"errors"
	"github.com/hauke96/sigolo"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Channels is a type that represents an m3u playlist containing 0 or more tracks
type Channels []Channel

// Track represents an m3u track
type Channel struct {
	Name      string
	Length    int
	URI       string
	TvgID     string
	TvgLogo   string
	TvgName   string
	GroupName string
}

var (
	tvgNameRegex    = regexp.MustCompile(`tvg-name="([^"]+)"`)
	tvgLogoRegex    = regexp.MustCompile(`tvg-logo="([^"]+)"`)
	tvgIDRegex      = regexp.MustCompile(`tvg-id="([^"]+)"`)
	groupTitleRegex = regexp.MustCompile(`group-title="([^"]+)"`)
	space           = regexp.MustCompile(`\s+`)
)

// Parse parses an m3u playlist with the given file name and returns a Channels
func Parse(fileName string) (channels Channels, err error) {
	f, err := os.Open(fileName)
	if err != nil {
		sigolo.Fatal("file not found")
		return
	}
	defer f.Close()

	var (
		onFirstLine = true
		scanner     = bufio.NewScanner(f)
	)

	for scanner.Scan() {
		line := scanner.Text()
		if onFirstLine && !strings.HasPrefix(line, "#EXTM3U") {
			err = errors.New("invalid m3u file format. Expected #EXTM3U file header")
			return
		}

		onFirstLine = false

		if strings.HasPrefix(line, "#EXTINF") {
			line := strings.Replace(line, "#EXTINF:", "", -1)
			// At this point the line will be something like "1 xxxxxxx"
			// We need "1, xxxxxx"
			tempInfo := strings.Split(line, " ")
			tempLength := tempInfo[0] // This is "1"
			if !strings.HasSuffix(tempLength, ",") {
				// We don't have a comma so we need to add it
				line = line[len(tempLength):]
				line = tempLength + ", " + line
			}
			trackInfo := strings.Split(line, ",")
			if len(trackInfo) < 2 {
				err = errors.New("invalid m3u file format. Expected EXTINF metadata to contain track length and name data")
				return
			}
			length, parseErr := strconv.Atoi(trackInfo[0])
			if parseErr != nil {
				err = errors.New("unable to parse length. Line: " + line)
				return
			}

			var (
				trackName = strings.Join(trackInfo[1:], " ")
				tvgName   string
				tvgID     string
				tvgLogo   string
				GroupName string
			)

			tvgLogoFind := tvgLogoRegex.FindStringSubmatch(trackName)
			if len(tvgLogo) != 0 {
				tvgLogo = tvgLogoFind[0]
				tvgLogo = strings.Replace(GroupName, "tvg-logo=\"", "", -1)
				tvgLogo = strings.Replace(GroupName, "\"", "", -1)
			}

			groupFind := groupTitleRegex.FindStringSubmatch(trackName)
			if len(groupFind) != 0 {
				GroupName = groupFind[0]
				GroupName = strings.Replace(GroupName, "group-title=\"", "", -1)
				GroupName = strings.Replace(GroupName, "\"", "", -1)
			}

			nameFind := tvgNameRegex.FindStringSubmatch(trackName)
			if len(nameFind) != 0 {
				tvgName = strings.Replace(nameFind[0], "tvg-name=\"", "", -1)
				tvgName = space.ReplaceAllString(tvgName, " ")
				tvgName = strings.Replace(tvgName, "\"", "", -1)
				tvgName = strings.TrimSpace(tvgName)
			}

			idFind := tvgIDRegex.FindStringSubmatch(trackName)
			if len(idFind) != 0 {
				tvgID = idFind[0]
				tvgID = strings.Replace(tvgID, "tvg-id=\"", "", -1)
				tvgID = strings.Replace(tvgID, "\"", "", -1)
			}

			track := &Channel{trackName, length, "", tvgID, tvgLogo, tvgName, GroupName}
			channels = append(channels, *track)
		} else if strings.HasPrefix(line, "#") || line == "" {
			continue
		} else if len(channels) == 0 {
			err = errors.New("URI provided for playlist with no tracks")
			return
		} else {
			channels[len(channels)-1].URI = line
		}
	}

	return channels, nil
}
