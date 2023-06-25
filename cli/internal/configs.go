package internal

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/urfave/cli/v2"
)

func GenConfigCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name: "config",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "app",
					Required: true,
					Usage:    "Specify an app to run your config command on",
				},
				&cli.StringFlag{
					Name:  "environment",
					Usage: "Specify an environment to run your apps:list command on",
				},
			},
			Usage: "This command prints an App's Configuration variables.",
			Action: func(ctx *cli.Context) error {
				c := NewConfigF(ctx)
				if err := c.GetConfiguration(ctx); err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
	}
}

func (c *Config) GetConfiguration(ctx *cli.Context) error {
	var err error

	appId, err := c.getAppIDFromFlags(ctx)
	if err != nil {
		return err
	}

	app, err := c.client.GetApp(appId)
	if err != nil {
		return err
	}

	if app.Env != nil {
		data, err := json.Marshal(app.Env)
		if err != nil {
			return err
		}
		out := make(map[string]interface{})
		json.Unmarshal(data, &out)

		for key, value := range out {
			fmt.Printf("%s=%s\n", key, value)
		}
	}

	return nil
}

// get
