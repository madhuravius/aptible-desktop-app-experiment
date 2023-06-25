package internal

import (
	"fmt"
	"log"

	"github.com/urfave/cli/v2"
)

func GenLogDrainsCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name: "log_drain:list",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "environment",
					Usage: "Specify an environment to run your log_drain:list command on",
				},
			},
			Usage: "This command lists all Log Drains.",
			Action: func(ctx *cli.Context) error {
				c := NewConfigF(ctx)
				if err := c.ListLogDrains(ctx); err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
	}
}

func (c *Config) ListLogDrains(ctx *cli.Context) error {
	var err error

	envs, err := c.getEnvironmentsFromFlags(ctx)
	if err != nil {
		return err
	}

	for _, env := range envs {
		logDrains, err := c.client.GetLogDrains(env.ID)
		if err != nil {
			return err
		}
		if len(logDrains) == 0 {
			continue
		}

		fmt.Printf("=== %s\n", env.Handle)
		for _, logDrain := range logDrains {
			fmt.Println(logDrain.Handle)
		}
	}
	return nil
}
