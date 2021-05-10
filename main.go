package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	defaultParallel = 10

	httpPrefix         = "http"
	parallelFlagPrefix = "-parallel"
)

func main() {
	parallel := flag.Int("parallel", defaultParallel, "amount of parallel requests")
	flag.Parse()

	urls := extractUrls(os.Args[1:])
	result := getHashResponses(urls, *parallel)
	for u, hash := range result {
		fmt.Println(u, hash)
	}
}

func extractUrls(args []string) []string {
	var urls []string
	for _, a := range args {
		if a == parallelFlagPrefix {
			continue
		}

		// check if the argument is not a numeric
		if _, err := strconv.Atoi(a); err != nil {
			urls = append(urls, a)
		}
	}

	return urls
}

type httpResponse struct {
	url  string
	hash string
}

func getHashResponses(urls []string, parallel int) map[string]string {
	// ensuring that we will have at least one worker
	if parallel < 1 {
		parallel = 1
	}

	urls = removeDuplicateUrl(urls)

	hashes := make(map[string]string)
	httpResponses := make(chan httpResponse)
	processUrl := make(chan string)

	defer close(httpResponses)
	defer close(processUrl)

	startWorkers(processUrl, httpResponses, parallel)

	var wg sync.WaitGroup
	wg.Add(len(urls))
	go func() {
		for _, u := range urls {
			processUrl <- adjustUrl(u)
		}
	}()

	go func() {
		for r := range httpResponses {
			hashes[r.url] = r.hash
			wg.Done()
		}
	}()

	wg.Wait()

	return hashes
}

func removeDuplicateUrl(urls []string) []string {
	var newUrls []string
	addedUrl := make(map[string]bool)
	for _, u := range urls {
		if !addedUrl[u] {
			addedUrl[u] = true
			newUrls = append(newUrls, u)
		}
	}

	return newUrls
}

func startWorkers(urls chan string, httpResponses chan httpResponse, parallel int) {
	for i := 0; i < parallel; i++ {
		go func(urls chan string, httpResponses chan httpResponse) {
			for u := range urls {
				r := httpResponse{url: u}
				b, err := doRequest(u)
				if err != nil {
					httpResponses <- r
					continue
				}
				r.hash = applyHash(b)
				httpResponses <- r
			}
		}(urls, httpResponses)
	}
}

func doRequest(url string) ([]byte, error) {
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func applyHash(b []byte) string {
	hash := md5.New()
	hash.Write(b)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func adjustUrl(url string) string {
	if strings.HasPrefix(url, httpPrefix) {
		return url
	}

	return fmt.Sprintf("%s://%s", httpPrefix, url)
}
