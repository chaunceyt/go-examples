package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/golang/glog"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "rctl"
	app.Usage = "Rekognition Command line tool"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		{
			Name:  "get-labels",
			Usage: "Detects instances of real-world entities within an image.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "source",
					Usage: "Name of the file to process.",
				},
			},
			Action: func(c *cli.Context) error {

				source := c.String("source")

				sess := session.New(&aws.Config{
					Region: aws.String("us-east-1"),
				})
				svc := rekognition.New(sess)

				labelInput := &rekognition.DetectLabelsInput{
					Image: &rekognition.Image{
						Bytes: getImage(source),
					},
				}
				labelResults, errResults := svc.DetectLabels(labelInput)
				if errResults != nil {
					fmt.Println("Error Detecting objects in image.")
					fmt.Println(errResults.Error())
					os.Exit(1)
				}

				fmt.Println("Number of labels found: ", len(labelResults.Labels))
				fmt.Println("")

				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Name", "Confidence"})

				for i := 0; i < len(labelResults.Labels); i++ {
					fmt.Println(*labelResults.Labels[i].Name, *labelResults.Labels[i].Confidence)
				}

				fmt.Println("")

				textInput := &rekognition.DetectTextInput{
					Image: &rekognition.Image{
						Bytes: getImage(source),
					},
				}
				textResults, errResults := svc.DetectText(textInput)
				if errResults != nil {
					fmt.Println("Error Detecting text in image.")
					fmt.Println(errResults.Error())
					os.Exit(1)
				}

				if len(textResults.TextDetections) > 0 {
					fmt.Println("Text found within the image")
					for j := 0; j < len(textResults.TextDetections); j++ {
						fmt.Println(*textResults.TextDetections[j].DetectedText)
					}
				}

				return nil
			},
		},
		{
			Name:  "detect-faces",
			Usage: "Detects faces within an image.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "source",
					Usage: "Name of the file to process.",
				},
			},
			Action: func(c *cli.Context) error {
				source := c.String("source")

				sess := session.New(&aws.Config{
					Region: aws.String("us-east-1"),
				})
				svc := rekognition.New(sess)

				input := &rekognition.DetectFacesInput{
					Image: &rekognition.Image{
						Bytes: getImage(source),
					},
					Attributes: []*string{aws.String("ALL")},
				}

				result, err := svc.DetectFaces(input)
				if err != nil {
					fmt.Println("Error Detecting objects in image.")
					fmt.Println(err.Error())
					os.Exit(1)
				}

				//table := tablewriter.NewWriter(os.Stdout)

				fmt.Println("Total faces detected", len(result.FaceDetails))
				fmt.Println("")
				for i := 0; i < len(result.FaceDetails); i++ {
					fmt.Println(*result.FaceDetails[i].Gender.Value, *result.FaceDetails[i].Gender.Confidence)
					fmt.Println("Age range", *result.FaceDetails[i].AgeRange.Low, "-", *result.FaceDetails[i].AgeRange.High, *result.FaceDetails[i].Confidence)
					fmt.Println(" > Emotions")
					//confidenceGender := strconv.FormatFloat(*result.FaceDetails[i].Gender.Confidence, 'E', -1, 64)

					//table.Append([]string{*result.FaceDetails[i].Gender.Value, confidenceGender})

					for j := 0; j < len(result.FaceDetails[i].Emotions); j++ {
						//confidence := strconv.FormatFloat(*result.FaceDetails[i].Emotions[j].Confidence, 'E', -1, 64)
						//table.Append([]string{*result.FaceDetails[i].Emotions[j].Type, confidence})
						fmt.Println(*result.FaceDetails[i].Emotions[j].Type, *result.FaceDetails[i].Emotions[j].Confidence)
					}
					//table.Append([]string{"", ""})
					fmt.Println("")
					//table.SetAutoWrapText(false)
					//table.Render()
				}

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// get Image - source can be the name of a file or a url
// @TODO: confirm we have an actual image
// https://socketloop.com/tutorials/golang-how-to-verify-uploaded-file-is-image-or-allowed-file-types
func getImage(source string) []byte {
	_, err := url.ParseRequestURI(source)
	if err != nil {
		file, err := os.Open(source)
		if err != nil {
			glog.Error(err)
		}

		defer file.Close()
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			glog.Error(err)
		}
		return bytes

	} else {
		resp, err := http.Get(source)
		if err != nil {
			glog.Error(err)
		}
		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			glog.Error(err)
		}
		return bytes

	}

}
