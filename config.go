package main

import (
	"github.com/hauke96/sigolo"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Target    string
	publicURL string
	epgURL    string
	Include   []string
	CacheTime int
}

func LoadConfig() (*Config, error) {
	target := os.Getenv("TARGET")
	if target == "" {
		sigolo.Fatal("env TARGET is required")
		return nil, nil
	}

	publicURL := os.Getenv("PUBLIC_URL")
	if publicURL == "" {
		sigolo.Fatal("env PUBLIC_URL is required")
		return nil, nil
	}

	epgURL := os.Getenv("EPG_URL")
	if publicURL == "" {
		sigolo.Fatal("env EPG_URL is required")
		return nil, nil
	}

	includeCategoriesRaw := os.Getenv("INCLUDE")
	if includeCategoriesRaw == "" {
		sigolo.Fatal("env INCLUDE is required")
		return nil, nil
	}

	includeCategories := strings.Split(includeCategoriesRaw, ",")

	cacheTimeString := os.Getenv("CACHE_TIME")
	if cacheTimeString == "" {
		sigolo.Fatal("env CACHE_TIME is required")
		return nil, nil
	}

	cacheTime, err := strconv.Atoi(cacheTimeString)
	if err != nil {
		return nil, err
	}

	if cacheTime < 1 {
		cacheTime = 1
	}

	return &Config{
		Target:    target,
		publicURL: publicURL,
		epgURL:    epgURL,
		Include:   includeCategories,
		CacheTime: cacheTime,
	}, nil
}
