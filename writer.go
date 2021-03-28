package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// writeLines writes the lines to the given file.
func writeLines(conf *Config, channels Channels, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	_, _ = fmt.Fprintln(w, `#EXTM3U`)

	for _, line := range channels {
		var icon = line.TvgLogo
		if icon == "" {
			icon = strings.Replace(icon, "Sd", "", -1)
			icon = strings.Replace(icon, "SD", "", -1)
			icon = strings.TrimSpace(icon)

			icon = strings.Replace(line.TvgName, "  ", " ", -1)
			icon = strings.Replace(icon, " ", "_", -1)
			icon = strings.Replace(icon, "FHD", "HD", -1)

			icon = fmt.Sprintf("http://%s/icons/%s.png", conf.publicURL, icon)
		}

		fmt.Fprintf(w, `#EXTINF:%d tvg-id="%s" tvg-name="%s" tvg-logo="%s" group-title="%s",%s`, line.Length, line.TvgID, line.TvgName, icon, line.GroupName, line.TvgName)
		fmt.Fprintln(w, ``)
		fmt.Fprintln(w, line.URI)
	}
	return w.Flush()
}
