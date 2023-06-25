package internal

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

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
