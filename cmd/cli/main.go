package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cbebe/worktracker/pkg/work"
)

func printUsage() {
	fmt.Println("USAGE:", os.Args[0], "start|stop|get [type]")
}

func handleCommand(args []string, service work.WorkService) error {
	switch strings.ToLower(args[0]) {
	case "start":
		if len(args) >= 2 {
			return service.StartLog(args[1])
		}
		return service.StartWork()
	case "stop":
		if len(args) >= 2 {
			return service.StopLog(args[1])
		}
		return service.StopWork()
	case "get":
		if len(args) >= 2 {
			works, err := service.GetWorkType(args[1])
			if works != nil {
				work.PrintWork(os.Stdout, works)
			}
			return err
		}
		works, err := service.GetWork()
		if works != nil {
			work.PrintWork(os.Stdout, works)
		}
		return err
	default:
		printUsage()
		return nil
	}
}

func main() {
	service, err := work.NewWorkService("work.db")
	if err != nil {
		log.Fatal(err)
	}

	if args := os.Args[1:]; len(args) > 0 {
		if err := handleCommand(args, service); err != nil {
			log.Fatal(err)
		}
	} else {
		printUsage()
	}
}
