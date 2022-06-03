package main

import (
	"log"
	"os"
	"strconv"

	"github.com/cbebe/worktracker/pkg/web"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("port not set: %v", err)
	}
	web.RunServer(port, "work.db")
}
