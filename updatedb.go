package main

import (
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
)

func updateDB() {
	q := "insert into insta (author_name, thumbnail_url, title, html, author_url, author_id, code) values "

	// generate parmeters for each ioresps (need 7 params per insta)
	for i, _ := range ioresps {
		p1 := strconv.Itoa(1 + (i * 7))
		p2 := strconv.Itoa(2 + (i * 7))
		p3 := strconv.Itoa(3 + (i * 7))
		p4 := strconv.Itoa(4 + (i * 7))
		p5 := strconv.Itoa(5 + (i * 7))
		p6 := strconv.Itoa(6 + (i * 7))
		p7 := strconv.Itoa(7 + (i * 7))

		params := "($" + p1 + ",$" + p2 + ",$" + p3 + ",$" + p4 + ",$" + p5 + ",$" + p6 + ",$" + p7 + "),"
		q += params
		// remove final ,
	}

	// trim trailing comma
	query := TrimSuffix(q, ",")

	// on conflict do nothing
	query += " on conflict (code) do nothing"

	// slice of empty interfaces to take values to be passed into prepared statement
	vals := []interface{}{}

	// add relevant fields of ioresp to vals slice
	for _, value := range ioresps {
		vals = append(vals, value.AuthorName, value.ThumbnailURL, value.Title, value.HTML, value.AuthorURL, value.AuthorID, value.Code)
	}

	// connect to db
	db, err := sqlx.Connect("postgres", "user=travisjones dbname=pushsouth sslmode=disable")
	if err != nil {
		panic(err)
	}

	// prepare and execuate query
	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err)
	}

	res, err := stmt.Exec(vals...)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	stmt.Close()
	db.Close()
}
