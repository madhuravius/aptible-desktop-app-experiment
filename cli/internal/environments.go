package internal

import (
	"fmt"

	"github.com/aptible/go-deploy/aptible"
	"github.com/urfave/cli/v2"
)

func (c *Config) ListEnvironments(ctx *cli.Context) error {
	var envs []aptible.Environment
	var err error

	environmentId := ctx.Value("environment").(int64)
	if environmentId > 0 {
		environment, err := c.client.GetEnvironment(environmentId)
		if err != nil {
			return err
		}
		envs = []aptible.Environment{environment}
		fmt.Printf("Got environment successfully - %s (%d)\n", environment.Handle, environment.ID)
	} else {
		envs, err = c.client.GetEnvironments()
		if err != nil {
			return err
		}
	}

	for _, env := range envs {
		fmt.Println(env.Handle)
	}
	return nil
}
