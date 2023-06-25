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
				&cli.StringFlag{
					Name:  "environment",
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

	envs, err = c.getEnvironmentsFromFlags(ctx)
	if err != nil {
		return err
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
