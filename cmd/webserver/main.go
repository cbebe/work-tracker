package main

import (
	"fmt"
	"log"
	"net/http"

	worktracker "github.com/cbebe/worktracker"
)

const port = 8080

func main() {
	store, err := worktracker.NewSqliteWorkStore("work.db")

	if err != nil {
		log.Fatal(err)
	}

	server := worktracker.NewWorkServer(store)
	http.ListenAndServe(fmt.Sprintf(":%d", port), server)
}
