package internal

import (
	"fmt"
	"log"

	"github.com/urfave/cli/v2"
)

func GenDatabaseCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name: "db:list",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "environment",
					Usage: "Specify an environment to run your db:list command on",
				},
			},
			Usage: "This command lists Databases in an Environment.",
			Action: func(ctx *cli.Context) error {
				c := NewConfigF(ctx)
				if err := c.ListDatabases(ctx); err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
	}
}

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
