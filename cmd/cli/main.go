package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cbebe/worktracker"
)

func printUsage() {
	fmt.Printf("USAGE: %s start|stop\n", os.Args[0])
}

func main() {
	store, err := worktracker.NewSqliteWorkStore("work.db")

	if err != nil {
		log.Fatal(err)
	}
	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		return
	}
	switch strings.ToLower(args[0]) {
	case "start":
		store.StartWork()
	case "stop":
		store.StopWork()
	case "get":
		worktracker.PrintWork(os.Stdout, store.GetWork())
	default:
		printUsage()
	}
}
