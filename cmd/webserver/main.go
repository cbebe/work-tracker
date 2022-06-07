package main

import (
	"log"
	"os"
	"strconv"

	"github.com/cbebe/work-tracker/pkg/web"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	path := os.Getenv(("DB_PATH"))
	if err != nil {
		log.Fatalf("port not set: %v", err)
	}
	web.RunServer(port, path)
}
