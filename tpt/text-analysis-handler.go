package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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

	region := "us-east-1"
	profile := "default"

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(region)},
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
	}))

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

	// Detect KeyPhrase(s) - returns the key phrases or talking points and a confidence score to support that this is a key phrase.
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
