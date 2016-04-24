package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

// slice of InstaOembedReqs
var ioreqs []InstaOembedReq

// slice of InstaOembedResp
var ioresps []InstaOembedResp

// takes a URL and pointer to wg, scrapes insta page, and returns a slice of
// requests formatted for insta's open oembed API endpoint
func getIWR(url string, wg *sync.WaitGroup) {
	// GET
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	// convert dumped response body to string
	respString := string(body)

	// split the string
	shards := strings.Split(respString, "window._sharedData = ")
	instaJson := strings.Split(shards[1], ";</script>")

	// instaResp to hold json
	var iwr InstaWebResp

	// unmarshal json into &iwr
	if err := json.Unmarshal([]byte(instaJson[0]), &iwr); err != nil {
		panic(err)
	}

	// for each piece of media add shortcode and url and then append to slice
	for _, value := range iwr.EntryData.Profilepage[0].User.Media.Nodes {
		var ior InstaOembedReq
		ior.ShortCode = value.Code
		// omitting script -- make sure you pull it in on the page
		ior.Url = "https://api.instagram.com/oembed/?url=http://instagr.am/p/" +
			value.Code + "&omitscript=true"
		ioreqs = append(ioreqs, ior)
	}

	fmt.Println(iwr.EntryData.Profilepage[0].User.Username, "ready!")

	wg.Done()
}

func asyncGetIWR(urls []string) {
	var wg sync.WaitGroup
	for _, value := range urls {
		wg.Add(1)
		go getIWR(value, &wg)
		wg.Wait()
	}
}

// takes an InstaOembedReq and pointer to wg (for concurrency). hits insta
// oembed api endpoint and appends response to ioresps slice
func getOembed(ioreq InstaOembedReq, wg *sync.WaitGroup) {
	// new GET req using oembed URL
	resp, err := http.Get(ioreq.Url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// fresh IOResp to hold each response
	var ioresp InstaOembedResp

	// decode response into &ioresp
	if err := json.NewDecoder(resp.Body).Decode(&ioresp); err != nil {
		panic(err)
	}

	// add insta shortcode to ioresp data structure
	ioresp.Code = ioreq.ShortCode

	// append each response to the IOResps slice
	ioresps = append(ioresps, ioresp)
	fmt.Println(ioresp)

	wg.Done()
}

func asyncGetOembed() {
	var wg sync.WaitGroup
	for _, value := range ioreqs {
		wg.Add(1)
		go getOembed(value, &wg)
		wg.Wait()
	}
}

func main() {
	urls := []string{"http://instagram.com/thrashermag",
		"http://instagram.com/habitatskateboards",
		"http://instagram.com/shaqueefaog",
		"http://instagram.com/theboardr"}

	t0 := time.Now() // temporary benchmarking
	asyncGetIWR(urls)
	t1 := time.Now()
	asyncGetOembed()
	t2 := time.Now()
	fmt.Println("Finished getIWR: ", t1.Sub(t0))
	fmt.Println("Finished getOembed: ", t2.Sub(t1))
	fmt.Println(len(ioresps))
}
