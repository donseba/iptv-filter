package main

import (
	"fmt"
	"github.com/hauke96/sigolo"
	"github.com/jasonlvhit/gocron"
	"log"
	"net/http"
	"path"
)

var conf *Config

func executeCronJob() {
	err := processM3u8(conf)
	if err != nil {
		sigolo.Info(err.Error())
	}
	gocron.Every(uint64(conf.CacheTime)).Minute().Do(processM3u8)
	<-gocron.Start()
}

func main() {
	var err error
	conf, err = LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}
	sigolo.Info("starting up")
	go executeCronJob()

	http.HandleFunc("/iptv.m3u", ServeM3u)
	http.Handle("/icons/", http.StripPrefix("/icons/", http.FileServer(http.Dir("icons"))))

	err = http.ListenAndServe(":65341", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func ServeM3u(w http.ResponseWriter, r *http.Request) {
	fmt.Println("serving file")

	fp := path.Join("iptv.m3u")

	http.ServeFile(w, r, fp)
}
