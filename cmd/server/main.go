package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
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

	fmt.Printf("%s\n", b)
}
