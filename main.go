package main

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// slice of InstaOembedReqs
var ioreqs []InstaOembedReq

// slice of InstaOembedResp
var ioresps []InstaOembedResp

func main() {
	// connect to db
	db, err := sqlx.Connect("postgres", "user=travisjones dbname=pushsouth sslmode=disable")
	if err != nil {
		panic(err)
	}

	// select urls from db and drop them in urls slice
	var urls []string
	db.Select(&urls, "select url from author")
	db.Close()

	t0 := time.Now() // temporary benchmarking
	asyncGetIWR(urls)
	t1 := time.Now()
	asyncGetOembed()
	t2 := time.Now()
	fmt.Println("Finished getIWR: ", t1.Sub(t0))
	fmt.Println("Finished getOembed: ", t2.Sub(t1))
	fmt.Println(len(ioresps))
	updateDB()
	t3 := time.Now()
	fmt.Println("Finished DB update: ", t3.Sub(t2))
}
