package main

import (
	"fmt"
	"html/template"
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
	store, err := worktracker.NewStore(worktracker.GetPath(os.Stdout))
	if err != nil {
		log.Fatalf("error creating work store: %v\n", err)
	}
	s := worktracker.NewWorkService(store)

	server := worktracker.NewWorkHandler(os.Stdout, s, template.Must(template.ParseFiles("layout.html")))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), server))
}
