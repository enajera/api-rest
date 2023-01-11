package server

import (
	"log"
	"net/http"
	"time"
	"github.com/go-chi/chi"
	"github.com/spf13/viper"
)

type Server struct {
	server *http.Server
}

func NewServer(mux *chi.Mux) *Server {
	port := viper.GetString("port")
	s := &http.Server{
		Addr:           ":"+ port,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return &Server{s}
}

func (s *Server) Run(){
	log.Fatal(s.server.ListenAndServe())
}