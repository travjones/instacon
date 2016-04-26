package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

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
	// ioresp.RawJson = string(body)

	// append each response to the IOResps slice
	ioresps = append(ioresps, ioresp)
	// fmt.Println(ioresp)

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
