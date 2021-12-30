package worktracker

import (
	"net/http"
)

type WorkServer struct {
	store WorkStore
	http.Handler
}

func NewWorkServer(store WorkStore) *WorkServer {
	w := new(WorkServer)

	w.store = store

	router := http.NewServeMux()
	router.Handle("/start", http.HandlerFunc(w.startWorkHandler))
	router.Handle("/stop", http.HandlerFunc(w.stopWorkHandler))
	w.Handler = router

	return w
}

func (s *WorkServer) startWorkHandler(w http.ResponseWriter, r *http.Request) {
	s.store.StartWork()
	w.WriteHeader(http.StatusOK)
}

func (s *WorkServer) stopWorkHandler(w http.ResponseWriter, r *http.Request) {
	s.store.StopWork()
	w.WriteHeader(http.StatusOK)
}
