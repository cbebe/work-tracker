package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cbebe/work-tracker/pkg/work"
)

func RunServer(port int, path string) {
	service, err := work.NewWorkService(path)

	if err != nil {
		log.Fatal(err)
	}

	server := NewWorkController(service)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), server))
}
