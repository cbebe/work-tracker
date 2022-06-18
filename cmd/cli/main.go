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
	args := os.Args[1:]
	if len(args) <= 0 {
		worktracker.PrintUsage()
	} else if err := worktracker.HandleCommand(args, s); err != nil {
		log.Fatal(err)
	}
}
