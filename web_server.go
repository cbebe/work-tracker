package worktracker

import (
	"fmt"
	"io"
	"net/http"
)

type workHandler struct {
	output  io.Writer
	service webService
	layout  PageLayout
	http.Handler
}

type PageLayout interface {
	Execute(wr io.Writer, data any) error
}

type webService interface {
	StartWork() error
	StopWork() error
	Store
}

func NewWorkHandler(output io.Writer, service webService, layout PageLayout) *workHandler {
	s := workHandler{
		output:  output,
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

func (s *workHandler) handleError(w http.ResponseWriter, err error) {
	fmt.Fprint(s.output, err)
	w.WriteHeader(http.StatusInternalServerError)
}

func (s *workHandler) startWorkHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.service.StartWork(); err != nil {
		s.handleError(w, err)
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
		s.handleError(w, err)
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
		s.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	PrintWorks(w, works)
}

func (s workHandler) stopWorkHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.service.StopWork(); err != nil {
		s.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
