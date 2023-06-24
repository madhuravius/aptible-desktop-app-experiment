package internal

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func (c *Config) ListMetricDrains(ctx *cli.Context) error {
	var err error

	envs, err := c.getEnvironmentsFromFlags(ctx)
	if err != nil {
		return err
	}

	for _, env := range envs {
		metricDrains, err := c.client.GetMetricDrains(env.ID)
		if err != nil {
			return err
		}
		if len(metricDrains) == 0 {
			continue
		}

		fmt.Printf("=== %s\n", env.Handle)
		for _, metricDrain := range metricDrains {
			fmt.Println(metricDrain.Handle)
		}
	}
	return nil
}
