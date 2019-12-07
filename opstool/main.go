package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "Operations ToolKit"
	app.Usage = "Execute various ops commands"
	app.Version = "0.0.1"

	netAppFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Value: "cthorn.com",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "aws",
			Usage: "Manage AWS resources",
			Subcommands: cli.Commands{
				{
					Name:  "whoami",
					Usage: "Get Caller Identity",
					Action: func(c *cli.Context) error {
						whoami()
						return nil
					},
				},
				{
					Name:  "running-instances",
					Usage: "Get Running EC2 Instances",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name: "list-instances",
						},
						cli.StringFlag{
							Name:  "region",
							Value: "us-east-1",
						},
						cli.StringFlag{
							Name:  "profile",
							Value: "default",
						},
					},
					Action: func(c *cli.Context) error {
						runningInstances(c)
						return nil
					},
				},
				{
					Name:  "ebs-volumes",
					Usage: "Get list of EBS volumes",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "region",
							Value: "us-east-1",
						},
						cli.StringFlag{
							Name:  "profile",
							Value: "default",
						},
					},
					Action: func(c *cli.Context) error {
						ebsVolumes(c)
						return nil
					},
				},
		},
		{
			Name:  "net",
			Usage: "Network debugging commands",
			Subcommands: cli.Commands{
				{
					Name:  "ns",
					Usage: "Looks up the name servers for a particular Host",
					Flags: netAppFlags,
					Action: func(c *cli.Context) error {
						ns, err := net.LookupNS(c.String("host"))
						if err != nil {
							return err
						}
						for i := 0; i < len(ns); i++ {
							fmt.Println(ns[i].Host)
						}
						return nil
					},
				},
				{
					Name:  "ip",
					Usage: "Looks up the IP address for a particular host",
					Flags: netAppFlags,
					Action: func(c *cli.Context) error {
						ip, err := net.LookupIP(c.String("host"))
						if err != nil {
							fmt.Println(err)
							return err
						}
						for i := 0; i < len(ip); i++ {
							fmt.Println(ip[i])
						}
						return nil
					},
				},
				{
					Name:  "txt",
					Usage: "Looks up the TXT records for a particular host",
					Flags: netAppFlags,
					Action: func(c *cli.Context) error {
						txt, err := net.LookupTXT(c.String("host"))
						if err != nil {
							fmt.Println(err)
							return err
						}
						for i := 0; i < len(txt); i++ {
							fmt.Println(txt[i])
						}
						return nil
					},
				},
				{
					Name:  "cname",
					Usage: "Looks up the cname for a particular host",
					Flags: netAppFlags,
					Action: func(c *cli.Context) error {
						cname, err := net.LookupCNAME(c.String("host"))
						if err != nil {
							fmt.Println(err)
							return err
						}
						fmt.Println(cname)
						return nil
					},
				},
				{
					Name:  "mx",
					Usage: "Looks up the mx for a particular host",
					Flags: netAppFlags,
					Action: func(c *cli.Context) error {
						mx, err := net.LookupMX(c.String("host"))
						if err != nil {
							fmt.Println(err)
							return err
						}
						for i := 0; i < len(mx); i++ {
							fmt.Println(mx[i])
						}
						return nil
					},
				},
			},
		},
		{
			Name:  "curtime",
			Usage: "Display current time",
			Action: func(c *cli.Context) error {
				now := time.Now()
				prettyTime := now.Format(time.RubyDate)
				fmt.Println("The current time is ", prettyTime)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
