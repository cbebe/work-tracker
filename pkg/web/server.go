package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cbebe/worktracker/pkg/work"
)

func RunServer(port int, path string) {
	store, err := work.NewSqliteWorkStore(path)

	if err != nil {
		log.Fatal(err)
	}

	server := NewWorkServer(store)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), server))
}
