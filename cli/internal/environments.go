package internal

import (
	"fmt"
	"log"

	"github.com/aptible/go-deploy/aptible"
	"github.com/urfave/cli/v2"
)

func GenEnvironmentsCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name: "environment:list",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "environment",
					Usage: "Specify an environment to run your metric_drain:list command on",
				},
			},
			Usage: "This command lists all Environments.",
			Action: func(ctx *cli.Context) error {
				c := NewConfigF(ctx)
				if err := c.ListEnvironments(ctx); err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
	}
}
func (c *Config) ListEnvironments(ctx *cli.Context) error {
	var envs []aptible.Environment
	var err error

	envs, err = c.getEnvironmentsFromFlags(ctx)
	if err != nil {
		return err
	}

	for _, env := range envs {
		fmt.Println(env.Handle)
	}
	return nil
}
