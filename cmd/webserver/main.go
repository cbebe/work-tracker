package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/cbebe/worktracker"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("port not set: %v", err)
	}
	store, err := worktracker.NewSqliteWorkStore("work.db")

	if err != nil {
		log.Fatal(err)
	}

	server := worktracker.NewWorkServer(store)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), server))
}
