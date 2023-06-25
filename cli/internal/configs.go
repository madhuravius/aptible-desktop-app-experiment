package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/urfave/cli/v2"
)

func GenConfigCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name: "config",
			Flags: []cli.Flag{
				&cli.Int64Flag{
					Name:     "app",
					Value:    0,
					Required: true,
					Usage:    "Specify an app to run your config command on",
				},
				&cli.Int64Flag{
					Name:  "environment",
					Value: 0,
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

	appId := ctx.Value("app").(int64)
	environmentId := ctx.Value("environment").(int64)

	app, err := c.client.GetApp(appId)
	if err != nil {
		return err
	}

	if environmentId != 0 {
		environment, err := c.client.GetEnvironment(environmentId)
		if err != nil {
			return err
		}

		if app.EnvironmentID != environment.ID {
			return errors.New("error - app's environment and environment param do not match")
		}
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
