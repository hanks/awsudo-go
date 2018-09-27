package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/hanks/awsudo-go/cmd"
	"github.com/hanks/awsudo-go/cmd/agent"
	c "github.com/hanks/awsudo-go/configs"
	"github.com/hanks/awsudo-go/utils"
	v "github.com/hanks/awsudo-go/version"
	"github.com/urfave/cli"
)

var configFlag = cli.StringFlag{
	Name:  "config, c",
	Value: "",
	Usage: "Load configuration from `FILE`",
}

func init() {
	// setup log file and format
	logFile, logErr := os.OpenFile(c.DefaultLogPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("Fail to find", c.DefaultLogPath, "awsudo start failed")
		os.Exit(1)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	if c.DEBUG {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	} else {
		log.SetFlags(log.Ldate | log.Ltime)
	}

	log.Printf("Log to file: %s", c.DefaultLogPath)
}

func main() {
	app := cli.NewApp()

	app.Name = "awsudo"
	app.UsageText = `	awsudo configure [--config|-c awsudo.toml]                   # create awsudo config to the default path
	awsudo start-agent [--config|-c awsudo.toml]                 # start agent server to store credentials with expiration
	awsudo stop-agent                                            # stop agent server in background
	awsudo [--config|-c awsudo.toml] awsRoleName [aws commands]  # run aws commands with awsudo
	`
	app.Usage = "Automated AWS API access using a SAML compliant identity provider"
	app.Version = v.Version

	app.Flags = []cli.Flag{
		configFlag,
	}

	app.Commands = []cli.Command{
		{
			Name:      "start-agent",
			Usage:     "Start agent server to store credentials with expiration",
			UsageText: "awsudo start-agent [--config|-c awsudo.toml]",
			Flags: []cli.Flag{
				configFlag,
			},
			Action: func(c *cli.Context) error {
				confPath := c.String("config")
				agent.RunAgentServer(confPath)

				return nil
			},
		},
		{
			Name:      "configure",
			Usage:     "Set config to default path or specified path from input",
			UsageText: "awsudo configure [--config|-c awsudo.toml]",
			Flags: []cli.Flag{
				configFlag,
			},
			Action: func(c *cli.Context) error {
				confPath := c.String("config")
				cmd.RunConfigure(confPath)

				return nil
			},
		},
		{
			Name:      "stop-agent",
			Usage:     "Stop agent server in background, and do some cleanup tasks",
			UsageText: "awsudo stop-agent [--config|-c awsudo.toml]",
			Flags: []cli.Flag{
				configFlag,
			},
			Action: func(c *cli.Context) error {
				confPath := c.String("config")
				cmd.Stop(confPath)

				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Action = func(c *cli.Context) error {
		roleName := c.Args().Get(0)
		confPath := c.String("config")

		agent.RunAgentClient(confPath, roleName)

		length := len(c.Args())
		if length > 1 {
			command := c.Args()[1:]
			log.Printf("Start to execute command: %s", command)
			if err := utils.ExecCommand(command); err != nil {
				log.Fatalf("Error is: %v", err)
			}
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
