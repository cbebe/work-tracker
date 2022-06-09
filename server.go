package worktracker

import (
	"fmt"
	"log"
	"net/http"
)

func RunServer(port int, path string) {
	service, err := NewWorkService(path)
	if err != nil {
		log.Fatal(err)
	}

	server := NewWorkController(service)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), server))
}
