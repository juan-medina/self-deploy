package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (p *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			switch r.URL.Path {
			case "/random":
				p.randomRequest(w)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (p *Server) randomRequest(w http.ResponseWriter) {
	number := rand.Int()
	_, _ = fmt.Fprint(w, number)
}
