package main

import "net/http"

type Server struct {
	port string
}

func NewServer(port string) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Listen() error {
	err := http.ListenAndServe(s.port, nil)
	if err != nil {
		return err
	}
	return nil
}
