package internal

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func (c *Config) ListDatabases(ctx *cli.Context) error {
	envs, err := c.getEnvironmentsFromFlags(ctx)
	if err != nil {
		return err
	}

	for _, env := range envs {
		dbs, err := c.client.GetDatabases(env.ID)
		if err != nil {
			return err
		}
		if len(dbs) == 0 {
			continue
		}

		fmt.Printf("=== %s\n", env.Handle)
		for _, db := range dbs {
			fmt.Println(db.Handle)
		}
	}
	return nil
}

// backup

// clone

// create

// deprovision

// dump

// execute

// reload

// rename

// replicate

// restart

// tunnel

// url

// versions
