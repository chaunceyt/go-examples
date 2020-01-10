package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/comprehend"
)

type InputData struct {
	Message string
}

// jsonAPIHandler send json response.
func jsonAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var fromLang string
	var fromText string
	var textStr string

	decoder := json.NewDecoder(r.Body)
	var inputData InputData
	err := decoder.Decode(&inputData)
	if err != nil {
		fmt.Println(err)
	}

	sess := appSession()
	svc := comprehend.New(sess)

	details := InputData{
		Message: inputData.Message,
	}

	// Detect dominant language and if not en translate.
	langCode := detectLanguage(details.Message, svc)

	if langCode != "en" {
		// Translate to english
		log.Println("Detected language code:", langCode, "Translating...")
		fromLang = langCode
		fromText = details.Message
		textStr = translateText(langCode, details.Message, sess)
	} else {
		textStr = details.Message
	}
	// Process textStr

	// Detect Sentiment - returns the overall sentiment of a text (Positive, Negative, Neutral, or Mixed).
	log.Println("Detecting Sentiment of content...")
	sentimentInput := &comprehend.DetectSentimentInput{
		LanguageCode: aws.String(langCode),
		Text:         aws.String(textStr),
	}

	sentimentResp, errResp := svc.DetectSentiment(sentimentInput)
	if errResp != nil {
		// TODO return an error to the screen instead of exiting.
		fmt.Println("Got error calling DetectSentiment")
		fmt.Println(errResp.Error())
	}

	// Detect KeyPhrase(s) - returns the key phrases or talking points and a confidence score.
	phrases := detectKeyPhrases(langCode, textStr, svc)

	// Detect Entity Recognition
	entities := detectEntities(langCode, textStr, svc)

	// Detect Syntax tokens.
	syntaxTokens := detectSyntax(langCode, textStr, svc)

	jsonObject, _ := json.Marshal(SentimentResultsDetails{
		Message:                textStr,
		FromLanguage:           fromLang,
		FromText:               fromText,
		Success:                true,
		Sentiment:              *sentimentResp.Sentiment,
		SentimentScoreMixed:    *sentimentResp.SentimentScore.Mixed,
		SentimentScoreNegative: *sentimentResp.SentimentScore.Negative,
		SentimentScoreNeutral:  *sentimentResp.SentimentScore.Neutral,
		SentimentScorePositive: *sentimentResp.SentimentScore.Positive,
		KeyPhrases:             phrases,
		Entities:               entities,
		SyntaxTokens:           syntaxTokens,
	})

	w.Write(jsonObject)

}
