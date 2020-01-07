# Text Processing Tool

Simple tool to analyze text.

## Installing


When one thinks about affording the least privileges. We recommend one creates an AWS account for this tool and create an inline policy that looks something like this. The account only needs API access (no console required).

NOTE: Download the .csv file so you can export them as environmental variables later.

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "polly:SynthesizeSpeech",
                "comprehend:DetectSentiment",
                "comprehend:DetectEntities",
                "comprehend:DetectDominantLanguage",
                "translate:TranslateText",
                "comprehend:DetectSyntax",
                "comprehend:DetectKeyPhrases",
                "polly:DescribeVoices"
            ],
            "Resource": "*"
        }
    ]
}
```

After creating an account as described above export AWS variables.

```
export AWS_ACCESS_KEY_ID=YOUR_ACCESS_KEY_ID
export AWS_SECRET_ACCESS_KEY=YOUR_SECRET_KEY
export AWS_REGION=YOUR_AWS_REGION
```


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

## TODO

1. text-to-speech endpoint (UI)
2. translate endpoint (translate one-to-many languages - i,e english -> es, pt, etc)
3. additional logging
4. AWS secrets via ENV (done)
5. Consistent error handling
6. Auth for json api


