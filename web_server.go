package worktracker

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type workHandler struct {
	service *WorkService
	layout  PageLayout
	http.Handler
}

type PageLayout interface {
	Execute(wr io.Writer, data any) error
}

func NewWorkHandler(service *WorkService, layout PageLayout) *workHandler {
	s := workHandler{
		service: service,
		layout:  layout,
	}

	router := http.NewServeMux()
	router.HandleFunc("/all", s.getWorkHandler)
	router.HandleFunc("/start", s.startWorkHandler)
	router.HandleFunc("/stop", s.stopWorkHandler)
	router.HandleFunc("/", s.sendAllWorkHandler)
	s.Handler = router

	return &s
}

func handleError(w http.ResponseWriter, err error) {
	fmt.Fprint(os.Stdout, err)
	w.WriteHeader(http.StatusInternalServerError)
}

func (s *workHandler) startWorkHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.service.StartWork(); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type workPageData struct {
	PageTitle string
	Works     []Work
}

func (s *workHandler) sendAllWorkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	works, err := s.service.GetWork(ID)
	if err != nil {
		handleError(w, err)
		return
	}
	data := workPageData{
		PageTitle: "All Work",
		Works:     works,
	}
	s.layout.Execute(w, data)
}

func (s *workHandler) getWorkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	works, err := s.service.GetWork(ID)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	PrintWorks(w, works)
}

func (s workHandler) stopWorkHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.service.StopWork(); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
