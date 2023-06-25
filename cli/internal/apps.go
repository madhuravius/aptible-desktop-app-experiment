package internal

import (
	"fmt"
	"github.com/aptible/go-deploy/aptible"
	"github.com/urfave/cli/v2"
	"log"
)

func GenAppCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name: "apps",
			Flags: []cli.Flag{
				&cli.Int64Flag{
					Name:  "environment",
					Value: 0,
					Usage: "Specify an environment to run your apps:list command on",
				},
			},
			Usage: "This command lists Apps in an Environment.",
			Action: func(ctx *cli.Context) error {
				c := NewConfigF(ctx)
				if err := c.ListApps(ctx); err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
	}
}

func (c *Config) ListApps(ctx *cli.Context) error {
	var envs []aptible.Environment
	var err error

	environmentId := ctx.Value("environment").(int64)

	if environmentId != 0 {
		environment, err := c.client.GetEnvironment(environmentId)
		if err != nil {
			return err
		}
		envs = []aptible.Environment{environment}
	} else {
		envs, err = c.client.GetEnvironments()
		if err != nil {
			return err
		}
	}

	for _, env := range envs {
		apps, err := c.client.GetApps(env.ID)
		if err != nil {
			return err
		}
		if len(apps) == 0 {
			continue
		}

		fmt.Printf("=== %s\n", env.Handle)
		for _, app := range apps {
			fmt.Println(app.Handle)
		}
	}
	return nil
}

// config

// deploy

// rebuild

// restart

// services
