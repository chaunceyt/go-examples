package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// https://aws.amazon.com/comprehend/features/
func main() {
	listenPort := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/text-analysis", textAnalysisHandler)
	router.HandleFunc("/text-to-speech", textToSpeechHandler)
	router.HandleFunc("/api/json", jsonAPIHandler).Methods(http.MethodPost)

	// http.HandleFunc("/", webformHandler)
	log.Fatal(http.ListenAndServe(":"+*listenPort, router))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	message := "Text Processing Tool"
	w.Write([]byte(message))
}
