package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dvaganov/fiboser/pkg/delivery/fiborest"
)

type (
	restServer struct {
		fiboCtrl fiborest.Controller
		srv      *http.Server
	}
)

func (s *restServer) Run(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.fiboCtrl.GetFibonacciRange)

	s.srv = &http.Server{
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         fmt.Sprintf(":%d", port),
	}

	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	return nil
}

func (s *restServer) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
