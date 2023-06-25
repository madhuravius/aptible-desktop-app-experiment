package internal

import (
	"fmt"
	"log"

	"github.com/urfave/cli/v2"
)

func GenMetricDrainsCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name: "metric_drain:list",
			Flags: []cli.Flag{
				&cli.Int64Flag{
					Name:  "environment",
					Value: 0,
					Usage: "Specify an environment to run your metric_drain:list command on",
				},
			},
			Usage: "This command lists all Metric Drains.",
			Action: func(ctx *cli.Context) error {
				c := NewConfigF(ctx)
				if err := c.ListMetricDrains(ctx); err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
	}
}

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
