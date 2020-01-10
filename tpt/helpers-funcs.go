package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
)

func region() string {
	region := os.Getenv("AWS_REGION")

	if region == "" {
		region = "us-east-1"
		log.Println("Using", region, "as the default region")
	}

	return region

}

func appSession() *session.Session {
	region := region()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(region)},
		SharedConfigState: session.SharedConfigEnable,
	}))

	return sess
}

func detectSyntax(langCode string, textStr string, svc *comprehend.Comprehend) []SentimentSyntaxTokens {

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

	return syntaxTokens
}

func detectEntities(langCode string, textStr string, svc *comprehend.Comprehend) []SentimentEntities {
	// Detect Entity Recognition
	entitiesInput := &comprehend.DetectEntitiesInput{
		LanguageCode: aws.String(langCode),
		Text:         aws.String(textStr),
	}

	entitiesResp, errEntities := svc.DetectEntities(entitiesInput)
	if errEntities != nil {
		fmt.Println("Got error calling DetectEntities")
		fmt.Println(errEntities.Error())
	}

	var entities []SentimentEntities

	for _, e := range entitiesResp.Entities {
		entities = append(entities, SentimentEntities{
			Text:  *e.Text,
			Score: *e.Score,
			Type:  *e.Type,
		})
	}
	return entities
}

func detectKeyPhrases(langCode string, textStr string, svc *comprehend.Comprehend) []SentimentKeyPhrases {
	// Detect KeyPhrase(s) - returns the key phrases or talking points and a confidence score.
	log.Println("Detecting KeyPhrases in content...")
	keyPhraseInput := &comprehend.DetectKeyPhrasesInput{
		LanguageCode: aws.String(langCode),
		Text:         aws.String(textStr),
	}

	keyPhraseResp, errKeyPhrase := svc.DetectKeyPhrases(keyPhraseInput)
	if errKeyPhrase != nil {
		fmt.Println("Got error calling DetectEntities")
		fmt.Println(errKeyPhrase.Error())
	}

	var phrases []SentimentKeyPhrases

	for _, p := range keyPhraseResp.KeyPhrases {
		phrases = append(phrases, SentimentKeyPhrases{
			Text:  *p.Text,
			Score: *p.Score,
		})

	}
	return phrases
}

// secureHeaders - send secure headers
func secureHeaders(w http.ResponseWriter) {
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "same-origin")
	w.Header().Set("Vary", "Accept-Encoding")
	w.WriteHeader(http.StatusOK)
}
