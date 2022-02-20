package main

import (
	"fmt"
	"log"
	"time"
)

func CaptainsLog(db *DB, dateTag string) string {
	id, err := db.QueryCaptainsLog()
	if err != nil {
		log.Fatal(err.Error())
	}

	if id == "" {
		now := time.Now()
		return fmt.Sprintf("%s,%s/%s/%s", now.Format("2006-01-02"), dateTag, now.Format("2006"), now.Format("01"))
	} else {
		return id
	}
}
