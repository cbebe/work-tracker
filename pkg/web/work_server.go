package web

import (
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/cbebe/worktracker/pkg/work"
)

type WorkServer struct {
	service work.WorkService
	layout  *template.Template
	http.Handler
}

func NewWorkServer(store *work.SqliteWorkStore) *WorkServer {
	s := new(WorkServer)

	s.service = work.WorkService{SqliteWorkStore: store}
	s.layout = template.Must(template.ParseFiles("layout.html"))

	router := http.NewServeMux()
	router.Handle("/all", http.HandlerFunc(s.getWorkHandler))
	router.Handle("/start", http.HandlerFunc(s.startWorkHandler))
	router.Handle("/stop", http.HandlerFunc(s.stopWorkHandler))
	router.Handle("/", http.HandlerFunc(s.sendAllWorkHandler))
	s.Handler = router

	return s
}

func handleError(w http.ResponseWriter, err error) {
	fmt.Fprint(os.Stdout, err)
	w.WriteHeader(http.StatusInternalServerError)
}

func (s *WorkServer) startWorkHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.service.StartWork(); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type WorkPageData struct {
	PageTitle string
	Works     []work.Work
}

func (s *WorkServer) sendAllWorkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	works, err := s.service.GetWork()
	if err != nil {
		handleError(w, err)
		return
	}
	data := WorkPageData{
		PageTitle: "All Work",
		Works:     works,
	}
	s.layout.Execute(w, data)
}

func (s *WorkServer) getWorkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	works, err := s.service.GetWork()
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	work.PrintWork(w, works)
}

func (s *WorkServer) stopWorkHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.service.StopWork(); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
