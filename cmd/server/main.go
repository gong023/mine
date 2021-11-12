package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gong023/mine/internal/env"
	"github.com/gong023/mine/internal/handler"
	"github.com/gong023/mine/pkg/bybit"
)

func main() {
	srv := http.Server{
		Addr: ":8080",
	}

	http.Handle("/", http.HandlerFunc(mainHandler))

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err == nil {
			w.WriteHeader(http.StatusOK)
		} else {
			_, _ = fmt.Fprintln(os.Stderr, err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	cfg, err := env.New()
	if err != nil {
		return
	}
	cli := bybit.NewClient(cfg.BybitHost, cfg.BybitKey, cfg.BybitSec)
	h := handler.New(cfg, cli)
	decision, err := h.Start(r.Context(), b)
	if err != nil {
		return
	}
	fmt.Printf("%s\n", decision)
}
