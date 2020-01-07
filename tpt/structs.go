package main

// SentimentSyntaxTokens struct.
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
