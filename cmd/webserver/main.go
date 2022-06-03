package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/cbebe/worktracker/pkg/work"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("port not set: %v", err)
	}
	store, err := work.NewSqliteWorkStore("work.db")

	if err != nil {
		log.Fatal(err)
	}

	server := work.NewWorkServer(store)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), server))
}
