package main

import (
	"log"
	"os"

	"github.com/cbebe/worktracker"
)

func main() {
	store, err := worktracker.NewStore()
	if err != nil {
		log.Fatalf("error creating work store: %v\n", err)
	}
	s := worktracker.NewWorkService(store)
	if err := worktracker.HandleCommand(os.Stdout, os.Args, s); err != nil {
		log.Fatalln(err)
	}
}
