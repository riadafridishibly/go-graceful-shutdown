package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func reqLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("Method = %s, Path = %s, Took = %v",
			r.Method, r.URL.Path, time.Since(now))
	})
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)

	srv := &http.Server{
		Handler: reqLogMiddleware(mux),
		Addr:    ":8083",
	}

	go func() {
		<-sig
		log.Println("Graceful shutdown sequence initiated")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Println("Error shutting down server. err:", err)
		}
	}()

	log.Println("Server started at: http://localhost:8083/")
	if err := srv.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Http server stopped")
		} else {
			log.Fatal(err)
		}
	}
}
