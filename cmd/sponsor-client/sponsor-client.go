package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

var (
	assets     *string
	listenAddr *string
)

func init() {
	assets = flag.String("assets", "./static", "Path to the assets after \"ng build\"")
	listenAddr = flag.String("listen_addr", "localhost:9000", "Listening addr for the server")

	flag.Parse()
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(*assets)))
	mux.HandleFunc("/_ready", ready)
	mux.HandleFunc("/_healthz", health)

	srv := &http.Server{
		Addr:        *listenAddr,
		ReadTimeout: 5 * time.Second,
		Handler:     mux,
	}

	log.Fatal(srv.ListenAndServe())
}

func ready(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok\n"))
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok\n"))
}
