package internal

import (
	"fmt"
	"log"

	"github.com/urfave/cli/v2"
)

func (c *Config) ListDatabases(ctx *cli.Context) {
	envs, err := c.getEnvironmentsFromFlags(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, env := range envs {
		dbs, err := c.client.GetDatabases(env.ID)
		if err != nil {
			log.Fatal(err)
		}
		if len(dbs) == 0 {
			continue
		}

		fmt.Printf("=== %s\n", env.Handle)
		for _, db := range dbs {
			fmt.Println(db.Handle)
		}
	}
}
