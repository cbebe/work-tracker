package main

import (
	"log"
	"os"

	"github.com/cbebe/worktracker"
)

func main() {
	p := os.Getenv("DB_PATH")
	s, err := worktracker.NewWorkService(p)
	if err != nil {
		log.Fatal(err)
	}
	args := os.Args[1:]
	if len(args) <= 0 {
		worktracker.PrintUsage()
	} else if err := worktracker.HandleCommand(args, s); err != nil {
		log.Fatal(err)
	}
}
