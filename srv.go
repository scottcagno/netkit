package netkit

import (
	"net/http"
	"time"
	"log"
	"fmt"
)

type Server struct {
	*http.Server
}

func NewWebServer() *Server {
	return &Server{}
}

func (svr *Server) Serve(host string, handler http.Handler) {
	server := &http.Server{
		Addr: 			host,
		Handler: 		handler,
		ReadTimeout:	10*time.Second,
		WriteTimeout: 	10*time.Second,
		MaxHeaderBytes:	1 << 20,
		TLSConfig:		nil,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(fmt.Sprintf("TIMESTAMP: %v\nHmmz. Something went down.\n~Sever.", time.Now()))
	}
}