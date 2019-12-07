package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli"
)

func runningInstances(c *cli.Context) {
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
	resultRegions, err := svc.DescribeRegions(nil)
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println("> Checking all regions for running EC2 instances...")

	for _, region := range resultRegions.Regions {
		fmt.Println(*region.RegionName)
		sess, err := session.NewSessionWithOptions(session.Options{
			Profile: profile,
			Config: aws.Config{
				Region: aws.String(*region.RegionName),
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		svc := ec2.New(sess)

		// Setup filters to get running, pending instances.
		params := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String("instance-state-name"),
					Values: []*string{
						aws.String("running"),
						aws.String("pending"),
					},
				},
			},
		}
		result, err := svc.DescribeInstances(params)
		if err != nil {
			fmt.Println("Error", err)
		}

		fmt.Println(" > found", len(result.Reservations))

		if c.Bool("list-instances") {

			// Loop through reservations to get the instances.
			for _, reservation := range result.Reservations {
				for _, instance := range reservation.Instances {
					instanceID := *instance.InstanceId
					instanceType := *instance.InstanceType
					publicDNS := *instance.PublicDnsName
					launchTime := *instance.LaunchTime
					keyname := *instance.KeyName
					fmt.Println(publicDNS, instanceID, instanceType, launchTime, keyname)
					break
				}
			}
		}
	}
}
