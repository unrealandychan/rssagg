package main

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func urlToFeed(url string) (RSSFeed, error) {
	httpClient := http.Client{
		Timeout: time.Second * 10,
	}
	response, err := httpClient.Get(url)
	if err != nil {
		return RSSFeed{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(response.Body)

	dat, err := io.ReadAll(response.Body)
	if err != nil {
		return RSSFeed{}, err
	}

	resFeed := RSSFeed{}
	err = xml.Unmarshal(dat, &resFeed)
	if err != nil {
		return RSSFeed{}, err
	}
	return resFeed, nil
}
