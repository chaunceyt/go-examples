package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli"
)

func ebsVolumes(c *cli.Context) {
	profile := c.String("profile")
	region := c.String("region")

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: profile,
		Config: aws.Config{
			Region: aws.String(region),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	svc := ec2.New(sess)
	input := &ec2.DescribeVolumesInput{}

	result, err := svc.DescribeVolumes(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}

	//	fmt.Println(result)
	for _, volume := range result.Volumes {
		for _, a := range volume.Attachments {
			fmt.Println(*volume.VolumeId, *volume.AvailabilityZone, *volume.VolumeType, *volume.Size, *volume.State, *a.InstanceId, *a.Device)
		}
	}
}
