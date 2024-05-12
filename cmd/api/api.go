package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Api struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	port     int
}

func NewApiEndpoint(infoLog, errorLog *log.Logger) *Api {
	return &Api{
		infoLog:  infoLog,
		errorLog: errorLog,
		port:     5555,
	}
}

func (api *Api) Run(ctx context.Context) error {
	routes := api.routes()

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", api.port),
		Handler:           routes,
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	api.infoLog.Printf("Starting Readmodels Api Server on port %d", api.port)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			api.errorLog.Printf("Readmodels Api Server failed, error: %s\n", err)
		}
	}()

	<-ctx.Done()
	return server.Shutdown(ctx)
}
