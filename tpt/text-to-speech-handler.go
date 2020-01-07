package main

import "net/http"

func textToSpeechHandler(w http.ResponseWriter, r *http.Request) {
	message := "Text to Speech using AWS Polly"
	w.Write([]byte(message))
}
