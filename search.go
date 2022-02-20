package main

import "log"

func Search(db *DB, searchTitleOnly bool, searchTerm string) string {
	var results []Result
	var err error

	if searchTitleOnly {
		results, err = db.QueryTitles(searchTerm)
	} else {
		results, err = db.QueryText(searchTerm)
	}

	if err != nil {
		log.Fatal(err.Error())
	}

	return serialize(results)
}
