package main

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
