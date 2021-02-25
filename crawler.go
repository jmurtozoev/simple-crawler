package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"errors"
)

type Crawler interface {
	Crawl(url string) (string, error)
	visited(url string) (string, bool)
}

type crawler struct {
	crawled map[string]string
	mux     sync.Mutex
}

func newCrawler() Crawler {
	return &crawler{
		crawled: make(map[string]string),
	}
}

func (c *crawler) visited(url string) (string, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	v, ok := c.crawled[url]
	if ok {
		return v, true
	}

	return "", false
}

func (c *crawler) save(url, title string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.crawled[url] = title
}

// Crawl page and return title
func (c *crawler) Crawl(url string) (title string, err error) {
	title, ok := c.visited(url)
	if ok {
		return title, err
	}

	response, err := http.Get(url)
	if err != nil {
		return title, errors.New(fmt.Sprintf("fetching url error: %v", err))
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return title, errors.New("reading response body error")
	}

	pageContent := string(body)

	startIndex := strings.Index(pageContent, "<title>")
	if startIndex == -1 {
		return title, errors.New("no title element found")
	}

	// starting index of page title
	startIndex += len("<title>")

	// ending index of page title
	endIndex := strings.Index(pageContent, "</title>")
	if endIndex < 0 {
		return title, errors.New("no tag for closing title element found")
	}

	title = pageContent[startIndex:endIndex]

	return title, err
}
