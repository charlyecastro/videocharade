package main

import (
	"homework-charlyecastro/servers/summary/handlers"
	"log"
	"net/http"
	"os"
)

//main is the main entry point for the server
func main() {

	addr := os.Getenv("SUMMARYADDR")

	//if it's blank, default to ":80", which means
	//listen port 80 for requests addressed to any host
	if len(addr) == 0 {
		addr = ":80"
	}
	mux := http.NewServeMux()
	log.Printf("server is listening at %s...", addr)
	mux.HandleFunc("/v1/summary", handlers.SummaryHandler)
	log.Fatal(http.ListenAndServe(addr, mux))
}
