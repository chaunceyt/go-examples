package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/aws/aws-sdk-go/service/translate"
)

// https://aws.amazon.com/comprehend/features/

type SentimentSyntaxTokens struct {
	PartOfSpeechScore float64
	PartOfSpeechTag   string
	Text              string
}

// SentimentKeyPhrases struct.
type SentimentKeyPhrases struct {
	Text  string
	Score float64
}

// SentimentEntities struct.
type SentimentEntities struct {
	Text  string
	Type  string
	Score float64
}

// SentimentResultsDetails struct.
type SentimentResultsDetails struct {
	Message                string
	FromLanguage           string
	FromText               string
	Success                bool
	Sentiment              string
	SentimentScoreMixed    float64
	SentimentScoreNegative float64
	SentimentScoreNeutral  float64
	SentimentScorePositive float64
	KeyPhrases             []SentimentKeyPhrases
	Entities               []SentimentEntities
	SyntaxTokens           []SentimentSyntaxTokens
}

// Language struct.
type Language struct {
	LanguageCode string
	Score        float64
}

const webform = `
<!doctype html>
<html lang="en">
<head>
<!-- Required meta tags -->
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

<!-- Bootstrap CSS -->
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">

<title>Text Processor!</title>
</head>
<body>
<main role="main" class="container">
    {{if .Success }}
    <h1>The text was analyzed. See the results below.</h1>
    {{ if .FromLanguage }}
    <p>We had to translate the language from <strong>{{ .FromLanguage }}</strong> language code before we could process it.</p>
    <p>Original Message: {{ .FromText }}</p>
    {{ end }}
    <p>Message: {{ .Message }}</p>
        
    <h2>Sentiment: {{ .Sentiment }}</h2>
   
    <table class="table">
    <thead>
        <tr>
            <th scope="col">Sentiment</th>
            <th scope="col">Score</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>Mixed</td>
            <td>{{ .SentimentScoreMixed }}</td>
        </tr>
        <tr>
            <td>Negative</td>
            <td>{{ .SentimentScoreNegative }}</td>
        </tr>
        <tr>
            <td>Neutral</td>
            <td>{{ .SentimentScoreNeutral }}</td>
        </tr>
        <tr>
            <td>Positive</td>
            <td>{{ .SentimentScorePositive }}</td>
        </tr>
    </tbody>
    </table>

    <h2>Talking points found in text.</h2>
    <table class="table">
    <thead>
        <tr>
            <th scope="col">Key Phrase</th>
            <th scope="col">Score</th>
        </tr>
    </thead>
    <tbody>
    {{ range .KeyPhrases}}
        <tr>
            <td>{{.Text }}</td>
            <td>{{ .Score }}</td>
        </tr>
    {{ end }}
    </tbody>
    </table>

    <h2>Text syntax with parts of speech</h2>
    <table class="table">
    <thead>
        <tr>
            <th scope="col">Text</th>
            <th scope="col">Tag</th>
            <th scope="col">Score</th>
        </tr>
    </thead>
    <tbody>
    {{ range .SyntaxTokens }}
        <tr>
        <td>{{ .Text }}</td>
        <td>{{ .PartOfSpeechTag }}</td>
        <td>{{ .PartOfSpeechScore }}</td>
        </tr>
    {{ end }}
    </tbody>
    </table>


    <h2>Named Entities</h2>
    <table class="table">
    <thead>
        <tr>
            <th scope="col">Entitiy</th>
            <th scope="col">Category</th>
            <th scope="col">Score</th>
        </tr>
    </thead>
    <tbody>
    {{ range .Entities }}
        <tr>
            <td>{{ .Text }}</td>
            <td>{{ .Type }}</td>
            <td>{{ .Score }}</td>
        </tr>
    {{ end}}
    </tbody> 
    </table>

    {{else}}
    <h1>Text Processing Tool</h1>
    <form method="POST">
    <div class="form-group">
        <label>Enter text to be processed:</label><br />
        <textarea name="message" rows="10" cols="100"></textarea><br />
    </div>
    <button type="submit" class="btn btn-primary">Submit</button>
    </form>
   {{end}}
   </main>
   </body>
   </html>
    `

func main() {
	listenPort := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	http.HandleFunc("/", webformHandler)
	http.ListenAndServe(":"+*listenPort, nil)
}

func webformHandler(w http.ResponseWriter, r *http.Request) {
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

	renderjson := false
	if renderjson {
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

		fmt.Println(string(jsonObject))
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonObject)
	} else {
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

}

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
	req, resp := svc.TextRequest(&translate.TextInput{
		SourceLanguageCode: aws.String(langCode),
		TargetLanguageCode: aws.String(toLang),
		Text:               aws.String(textContent),
	})
	err := req.Send()
	if err != nil {
		return ""
	}
	return *resp.TranslatedText
}
