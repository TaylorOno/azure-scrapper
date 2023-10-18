package main

import (
	"azure-scrapper/internal/scrapper"
	"log"
	"net/http"
	"os"
)

func main() {
	listenAddr := ":9090"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}
	http.HandleFunc("/scrapper", scrapper.Handle)
	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
