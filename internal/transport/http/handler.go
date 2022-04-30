package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	Router  *mux.Router
	Service DatasourceService
	Server  *http.Server
}

func NewHandler(service DatasourceService) *Handler {
	h := &Handler{
		Service: service,
	}
	h.Router = mux.NewRouter()
	h.Router.Use(JSONMiddleware)
	h.Router.Use(LoggingMiddlerware)
	h.Router.Use(TimeoutMiddlerware)
	h.mapRoutes()

	h.Server = &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: h.Router,
	}

	return h
}

func (h *Handler) mapRoutes() {
	h.Router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})

	h.Router.HandleFunc("/api/v1/datasource", h.PostDatasource).Methods("POST")
	h.Router.HandleFunc("/api/v1/datasource/{id}", h.GetDatasource).Methods("GET")
	h.Router.HandleFunc("/api/v1/datasource/{id}", h.UpdateDatasource).Methods("PUT")
	h.Router.HandleFunc("/api/v1/datasource/{id}", h.DeleteDatasource).Methods("DELETE")
}

func (h *Handler) Serve() error {
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancle := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancle()
	h.Server.Shutdown(ctx)

	log.Println("shut down gracefully.")
	return nil
}
