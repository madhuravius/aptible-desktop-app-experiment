package internal

import (
	"fmt"
	"log"

	"github.com/aptible/go-deploy/aptible"
	"github.com/urfave/cli/v2"
)

func (c *Config) ListApps(ctx *cli.Context) {
	var envs []aptible.Environment
	var err error

	environmentId := ctx.Value("environment").(int64)

	if environmentId > 0 {
		environment, err := c.client.GetEnvironment(environmentId)
		if err != nil {
			log.Fatal(err)
		}
		envs = []aptible.Environment{environment}
	} else {
		envs, err = c.client.GetEnvironments()
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, env := range envs {
		apps, err := c.client.GetApps(env.ID)
		if err != nil {
			log.Fatal(err)
		}
		if len(apps) == 0 {
			continue
		}

		fmt.Printf("=== %s\n", env.Handle)
		for _, app := range apps {
			fmt.Println(app.Handle)
		}
	}
}
