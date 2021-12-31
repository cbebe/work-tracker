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

func handleCommand(cmd string, store worktracker.WorkStore) error {
	switch strings.ToLower(cmd) {
	case "start":
		return store.StartWork()
	case "stop":
		return store.StopWork()
	case "get":
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
		if err := handleCommand(args[0], store); err != nil {
			log.Fatal(err)
		}
	} else {
		printUsage()
	}
}
