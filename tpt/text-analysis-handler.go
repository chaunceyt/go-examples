package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/comprehend"
)

func textAnalysisHandler(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.New("webform").Parse(webform))

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	var fromLang string
	var fromText string
	var textStr string

	sess := appSession()

	svc := comprehend.New(sess)

	details := SentimentResultsDetails{
		Message: r.FormValue("message"),
	}

	// Detect dominant language and if not en translate.
	langCode := detectLanguage(details.Message, svc)

	if langCode != "en" {
		// Translate to english
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

	// Process template.
	tmpl.Execute(w, SentimentResultsDetails{
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

}
