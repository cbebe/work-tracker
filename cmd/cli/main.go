package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cbebe/worktracker"
)

func printUsage() {
	fmt.Printf("USAGE: %s start|stop|get [type]\n", os.Args[0])
}

func handleCommand(args []string, store worktracker.WorkStore) error {
	switch strings.ToLower(args[0]) {
	case "start":
		if len(args) >= 2 {
			return store.StartLog(args[1])
		}
		return store.StartWork()
	case "stop":
		if len(args) >= 2 {
			return store.StopLog(args[1])
		}
		return store.StopWork()
	case "get":
		if len(args) >= 2 {
			works, err := store.GetWorkType(args[1])
			if works != nil {
				worktracker.PrintWork(os.Stdout, works)
			}
			return err
		}
		works, err := store.GetWork()
		if works != nil {
			worktracker.PrintWork(os.Stdout, works)
		}
		return err
	default:
		printUsage()
		return nil
	}
}

func main() {
	store, err := worktracker.NewSqliteWorkStore("work.db")
	if err != nil {
		log.Fatal(err)
	}

	if args := os.Args[1:]; len(args) > 0 {
		if err := handleCommand(args, store); err != nil {
			log.Fatal(err)
		}
	} else {
		printUsage()
	}
}
