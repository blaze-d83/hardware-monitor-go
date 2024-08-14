package server

import (
	"fmt"
	"net/http"
	"sync"
)

type Server struct {
	subscriberMessageBuffer int
	subscribersMutex        sync.Mutex
	subscribers             map[*Subscriber]struct{}
	mux                     *http.ServeMux
}

func NewServer() *Server {
	s := &Server{
		subscriberMessageBuffer: 10,
		subscribers:             make(map[*Subscriber]struct{}),
		mux:                     http.NewServeMux(),
	}
	s.mux.Handle("/", http.FileServer(http.Dir("./static")))
	s.mux.HandleFunc("/ws", s.subscriberHandler)
	return s
}

func (s *Server) Start() error {
	fmt.Println("Starting server on port: 8080")
	return http.ListenAndServe(":8080", s.mux)
}
