package web

import (
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/cbebe/worktracker/pkg/work"
)

type WorkController struct {
	service work.WorkService
	layout  *template.Template
	http.Handler
}

func NewWorkController(service work.WorkService) *WorkController {
	s := new(WorkController)

	s.service = service
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

func (s *WorkController) startWorkHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *WorkController) sendAllWorkHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *WorkController) getWorkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	works, err := s.service.GetWork()
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	work.PrintWork(w, works)
}

func (s *WorkController) stopWorkHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.service.StopWork(); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
