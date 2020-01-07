package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// https://aws.amazon.com/comprehend/features/
func main() {
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	// Ensure the required environmental variables have been exported.
	if awsAccessKeyID == "" || awsSecretAccessKey == "" {
		errorMessage := "Please export the following\n $ export AWS_ACCESS_KEY_ID=YOUR_AKID\n $ export AWS_SECRET_ACCESS_KEY=YOUR_SECRET_KEY"
		log.Fatal(errorMessage)
	}

	listenPort := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/text-analysis", textAnalysisHandler)
	router.HandleFunc("/text-to-speech", textToSpeechHandler)
	router.HandleFunc("/api/json", jsonAPIHandler).Methods(http.MethodPost)
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		format := "%s - - [%s] \"%s %s %s\" %s\n"
		log.Printf(format, r.RemoteAddr, time.Now().Format(time.RFC1123),
			r.Method, r.URL.Path, r.Proto, r.UserAgent())
	})

	log.Println("Text Processing Service running on port:", *listenPort)
	log.Fatal(http.ListenAndServe(":"+*listenPort, router))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	message := "Text Processing Tool"
	w.Write([]byte(message))
}
