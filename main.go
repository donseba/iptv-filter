package main

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"log"
	"net/http"
	"path"
)

var conf *Config

func executeCronJob() {
	processM3u8(conf)

	gocron.Every(1).Hour().Do(processM3u8)
	<-gocron.Start()
}

func main() {
	var err error
	conf, err = LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	go executeCronJob()

	http.HandleFunc("/iptv.m3u", ServeM3u)
	http.Handle("/icons/", http.StripPrefix("/icons/", http.FileServer(http.Dir("icons"))))

	http.ListenAndServe(":65341", nil)
}

func ServeM3u(w http.ResponseWriter, r *http.Request) {
	fmt.Println("serving file")

	fp := path.Join("iptv.m3u")

	http.ServeFile(w, r, fp)
}
