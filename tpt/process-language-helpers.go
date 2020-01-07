package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/aws/aws-sdk-go/service/translate"
)

// detectLanguage - detect the language posted into the form.
func detectLanguage(textContent string, svc *comprehend.Comprehend) string {

	langDetectInput := &comprehend.DetectDominantLanguageInput{
		Text: aws.String(textContent),
	}

	log.Println("Detecting dominant language...")
	languageResp, errLanguageResp := svc.DetectDominantLanguage(langDetectInput)
	if errLanguageResp != nil {
		fmt.Println("Got an error calling DetectDominantLanguage")
		fmt.Println(errLanguageResp)
	}

	lr, err := json.Marshal(languageResp.Languages)
	if err != nil {
		fmt.Println(err)
	}
	var languages []Language
	json.Unmarshal(lr, &languages)
	langCode := languages[0].LanguageCode

	log.Println("Language detected", langCode)
	return langCode
}

// translateText - translate from the detected language to en
func translateText(langCode string, textContent string, sess *session.Session) string {
	toLang := "en"

	svc := translate.New(sess)
	log.Println("Translating the submitted text...")
	req, resp := svc.TextRequest(&translate.TextInput{
		SourceLanguageCode: aws.String(langCode),
		TargetLanguageCode: aws.String(toLang),
		Text:               aws.String(textContent),
	})
	err := req.Send()
	if err != nil {
		return ""
	}

	log.Println("Text translation completed...")
	return *resp.TranslatedText
}
