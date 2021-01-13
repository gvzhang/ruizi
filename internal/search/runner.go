package search

import (
	"net"
	"net/http"
	"ruizi/internal"
	"ruizi/internal/search/api"
	"time"
)

type Runner struct {
	server *http.Server
}

func NewRunner() *Runner {
	r := new(Runner)
	return r
}

func (r *Runner) Start() error {
	searchHttpServer := api.SearchHttpServer{}
	err := searchHttpServer.InitRouter()
	if err != nil {
		return err
	}

	var listener net.Listener
	addr := internal.GetConfig().Search.HttpListenAddr
	listener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	r.server = &http.Server{
		Handler:      searchHttpServer.Router,
		IdleTimeout:  2 * time.Hour,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	r.server.SetKeepAlivesEnabled(true)
	return r.server.Serve(listener)
}

func (r *Runner) Stop() error {
	return nil
}
