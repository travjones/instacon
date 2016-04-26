package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

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
	if err != nil {
		panic(err)
	}

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
