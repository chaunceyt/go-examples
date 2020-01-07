package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
)

type InputData struct {
	Message string
}

// jsonAPIHandler send json response.
func jsonAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var inputData InputData
	err := decoder.Decode(&inputData)
	if err != nil {
		fmt.Println(err)
	}

	var fromLang string
	var fromText string
	var textStr string

	region := os.Getenv("AWS_REGION")

	if region == "" {
		region = "us-east-1"
		log.Println("Using", region, "as the default region")
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(region)},
		SharedConfigState: session.SharedConfigEnable,
	}))

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
	log.Println("Detecting KeyPhrases in content...")
	keyPhraseInput := &comprehend.DetectKeyPhrasesInput{
		LanguageCode: aws.String(langCode),
		Text:         aws.String(textStr),
	}

	keyPhraseResp, errKeyPhrase := svc.DetectKeyPhrases(keyPhraseInput)
	if errKeyPhrase != nil {
		fmt.Println("Got error calling DetectEntities")
		fmt.Println(errResp.Error())
	}

	var phrases []SentimentKeyPhrases

	for _, p := range keyPhraseResp.KeyPhrases {
		phrases = append(phrases, SentimentKeyPhrases{
			Text:  *p.Text,
			Score: *p.Score,
		})

	}

	// Detect Entity Recognition
	entitiesInput := &comprehend.DetectEntitiesInput{
		LanguageCode: aws.String(langCode),
		Text:         aws.String(textStr),
	}

	entitiesResp, errEntities := svc.DetectEntities(entitiesInput)
	if errEntities != nil {
		fmt.Println("Got error calling DetectEntities")
		fmt.Println(errResp.Error())
	}

	var entities []SentimentEntities

	for _, e := range entitiesResp.Entities {
		entities = append(entities, SentimentEntities{
			Text:  *e.Text,
			Score: *e.Score,
			Type:  *e.Type,
		})
	}

	// Detect Syntax
	syntaxInput := &comprehend.DetectSyntaxInput{
		LanguageCode: aws.String(langCode),
		Text:         aws.String(textStr),
	}

	syntaxResp, errSyntax := svc.DetectSyntax(syntaxInput)
	if errSyntax != nil {
		fmt.Println("Got error calling DetectSyntax")
		fmt.Println(errSyntax)
	}

	var syntaxTokens []SentimentSyntaxTokens

	for _, st := range syntaxResp.SyntaxTokens {
		syntaxTokens = append(syntaxTokens, SentimentSyntaxTokens{
			PartOfSpeechScore: *st.PartOfSpeech.Score,
			PartOfSpeechTag:   *st.PartOfSpeech.Tag,
			Text:              *st.Text,
		})
	}

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
