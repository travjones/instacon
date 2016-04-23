package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
	"time"
)

// slice of InstaOembedResp
var ioresps []InstaOembedResp

// takes a URL, scrapes insta page, and returns a slice of requests formatted
// for insta's open oembed API endpoint
func getIWR(url string) []InstaOembedReq {
	// GET
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// dump response and include body
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}

	// convert dumped response body to string
	respString := string(dump)

	// split the string
	shards := strings.Split(respString, "window._sharedData = ")
	instaJson := strings.Split(shards[1], ";</script>")

	// instaResp to hold json
	var iwr InstaWebResp

	// unmarshal json into &iwr
	if err := json.Unmarshal([]byte(instaJson[0]), &iwr); err != nil {
		panic(err)
	}

	// slice of InstaOembedReqs
	var ioreqs []InstaOembedReq

	// for each piece of media add shortcode and url and then append to slice
	for i := 0; i < 12; i++ {
		var ior InstaOembedReq
		ior.ShortCode = iwr.EntryData.Profilepage[0].User.Media.Nodes[i].Code
		// omitting script -- make sure you pull it in on the page
		ior.Url = "https://api.instagram.com/oembed/?url=http://instagr.am/p/" +
			iwr.EntryData.Profilepage[0].User.Media.Nodes[i].Code +
			"&omitscript=true"
		ioreqs = append(ioreqs, ior)
	}

	return ioreqs
}

// takes an InstaOembedReq and pointer to wg (for concurrency). hits insta
// oembed api endpoint and appends response to ioresps slice
func getOembed(ior InstaOembedReq, wg *sync.WaitGroup) {
	// new GET req using oembed URL
	req, err := http.NewRequest("GET", ior.Url, nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	// Do request and drop response in resp
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close() // close dat shit or leak memory

	// fresh IOResp to hold each response
	var ioresp InstaOembedResp

	// decode response into &ioresp
	if err := json.NewDecoder(resp.Body).Decode(&ioresp); err != nil {
		panic(err)
	}

	// add shortcode to ioresp
	ioresp.ShortCode = ior.ShortCode

	// append each response to the IOResps slice
	ioresps = append(ioresps, ioresp)
	fmt.Println(ioresp)

	wg.Done()
}

func asyncGetOembed(urls []string) {
	for _, value := range urls {
		ioreqs := getIWR(value)
		var wg sync.WaitGroup
		for _, value := range ioreqs {
			wg.Add(1)
			go getOembed(value, &wg)
			wg.Wait()
		}
	}
}

func main() {
	urls := []string{"http://instagram.com/thrashermag",
		"http://instagram.com/habitatskateboards",
		"http://instagram.com/shaqueefaog",
		"http://instagram.com/theboardr"}

	t0 := time.Now()
	asyncGetOembed(urls)
	t1 := time.Now()
	fmt.Println(t1.Sub(t0))
	fmt.Println(len(ioresps))
}
