package server

import (
	"log"
	"net/http"
	"time"

	"github.com/St-Ivanov/todo-list-project/internal/api"
)

type HttpServer struct {
	Loger *log.Logger
	Serv  *http.Server
}

// A function that creates an http server with specific settings
func NewHttpServer(loger *log.Logger) *HttpServer {
	h := http.NewServeMux()

	api.Init(h)

	serv := &HttpServer{
		Loger: loger,
		Serv: &http.Server{
			Addr:         ":7540",
			Handler:      h,
			ErrorLog:     loger,
			ReadTimeout:  time.Duration(5) * time.Second,
			WriteTimeout: time.Duration(10) * time.Second,
			IdleTimeout:  time.Duration(15) * time.Second,
		},
	}

	return serv
}
