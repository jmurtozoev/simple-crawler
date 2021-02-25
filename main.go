package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"github.com/asaskevich/govalidator"
	"strings"
	"sync"
	"time"
)

type handler struct {
	c Crawler
}

var port = flag.Int("port", 8080, "server port")

func main() {
	flag.Parse()

	c := newCrawler()
	h := handler{c}

	http.HandleFunc("/crawl", h.GetTitles)

	log.Println("Running server on port:", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

type Url struct {
	url string
	title string
}

func (h handler) GetTitles(w http.ResponseWriter, r *http.Request){
	var wg sync.WaitGroup
	invalidUrls := make([]string, 0)
	urlTitles := make(map[string]string)
	ch := make(chan Url)


	urls, ok := r.URL.Query()["urls"]
	if !ok || len(urls) < 1 {
		http.Error(w, "url param 'urls' is missing", http.StatusBadRequest)
		return
	}

	for _, u := range urls {
		if valid(u) {
			wg.Add(1)
			go func(u string) {
				defer wg.Done()
				title, err := h.c.Crawl(u)
				if err != nil {
					log.Println("crawl error: ", err)
				}

				ch <- Url{url: u, title: title}
			}(u)

		} else {
			invalidUrls = append(invalidUrls, u)
		}
	}

	go func() {
		for u := range ch {
			urlTitles[u.url] = u.title
		}
	}()

	wg.Wait()
	time.Sleep(100*time.Microsecond)
	close(ch)

	resp := map[string]interface{}{
		"urls": urlTitles,
		"invalid_urls": invalidUrls,
	}

	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(resp)
	if err != nil {
		log.Fatal("marshalling response error: ", err)
	}

	if _, err = w.Write(b); err != nil {
		log.Fatal("writing response error: ", err)
	}
}

// validate url
func valid(u string) bool {
	if !strings.HasPrefix(u, "https://") {
		return false
	}

	return govalidator.IsURL(u)
}