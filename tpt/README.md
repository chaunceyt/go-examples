# Text Processing Tool

Simple tool to analyze text.

## Workflow

1. Detect language of text submitted using `comprehend.DetectDominantLanguage`
2. If not English translate to English using `translate.TextRequest` detected language to `en`
3. If English no need to translate
4. Get Sentiment Analysis using `comprehend.DetectSentiment` 
5. Get Key Phrases or talking points found in text using `comprehend.DetectKeyPhrases`
6. Get Entities using `comprehend.DetectKeyPhrases`
7. Get Syntax of text (Nouns, verbs, pronouns, etc) using `comprehend.DetectSyntax`

## Tech used.

1. AWS Comphrend AI Service
2. AWS Translate AI Service
3. Golang AWS SDK
4. AWS Account with `~/.aws/credentials` in place.

## TODO

1. text-to-speech endpoint (UI)
2. translate endpoint (translate one-to-many languages - i,e english -> es, pt, etc)
3. additional logging
4. AWS secrets via ENV
