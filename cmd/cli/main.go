package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cbebe/worktracker"
)

func printUsage() {
	fmt.Println("USAGE:", os.Args[0], "start|stop|get [type]")
}

func handleCommand(args []string, service *worktracker.WorkService) error {
	switch strings.ToLower(args[0]) {
	case "start":
		if len(args) >= 2 {
			return service.StartLog(args[1], worktracker.ID)
		}
		return service.StartWork()
	case "stop":
		if len(args) >= 2 {
			return service.StopLog(args[1], worktracker.ID)
		}
		return service.StopWork()
	case "get":
		if len(args) >= 2 {
			works, err := service.GetWorkType(args[1], worktracker.ID)
			if works != nil {
				worktracker.PrintWorks(os.Stdout, works)
			}
			return err
		}
		works, err := service.GetWork(worktracker.ID)
		if works != nil {
			worktracker.PrintWorks(os.Stdout, works)
		}
		return err
	default:
		printUsage()
		return nil
	}
}

func main() {
	path := os.Getenv("DB_PATH")
	service, err := worktracker.NewWorkService(path)
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
