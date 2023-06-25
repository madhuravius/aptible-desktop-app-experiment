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
		{
			Name: "db:backup",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "environment",
					Usage: "Specify an environment to run your db:backup command on",
				},
				&cli.StringFlag{
					Name:     "database",
					Required: true,
					Usage:    "Specify an app to run your db:backup command on",
				},
			},
			Usage: "This command is used to create Database Backups.",
			Action: func(ctx *cli.Context) error {
				c := NewConfigF(ctx)
				if err := c.BackupDatabase(ctx); err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
		{
			Name: "db:deprovision",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "environment",
					Usage: "Specify an environment to run your db:backup command on",
				},
				&cli.StringFlag{
					Name:     "database",
					Required: true,
					Usage:    "Specify an app to run your db:backup command on",
				},
			},
			Usage: "This command is used to deprovision a Database.",
			Action: func(ctx *cli.Context) error {
				c := NewConfigF(ctx)
				if err := c.DeprovisionDatabase(ctx); err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
		{
			Name: "db:reload",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "environment",
					Usage: "Specify an environment to run your db:backup command on",
				},
				&cli.StringFlag{
					Name:     "database",
					Required: true,
					Usage:    "Specify an app to run your db:backup command on",
				},
			},
			Usage: "This command reloads a Database by replacing the running Database Container with a new one.",
			Action: func(ctx *cli.Context) error {
				c := NewConfigF(ctx)
				if err := c.ReloadDatabase(ctx); err != nil {
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

func (c *Config) BackupDatabase(ctx *cli.Context) error {
	dbId, err := c.getDatabaseIDFromFlags(ctx)
	if err != nil {
		return err
	}

	op, err := c.client.DatabaseOperation(dbId, "backup")
	if err != nil {
		return err
	}

	if err = c.attachToOperationLogs(op); err != nil {
		return err
	}

	return nil
}

// clone

// create

func (c *Config) DeprovisionDatabase(ctx *cli.Context) error {
	dbId, err := c.getDatabaseIDFromFlags(ctx)
	if err != nil {
		return err
	}

	op, err := c.client.DatabaseOperation(dbId, "deprovision")
	if err != nil {
		return err
	}

	if err = c.attachToOperationLogs(op); err != nil {
		return err
	}

	return nil
}

// dump

// execute

// reload

func (c *Config) ReloadDatabase(ctx *cli.Context) error {
	dbId, err := c.getDatabaseIDFromFlags(ctx)
	if err != nil {
		return err
	}

	op, err := c.client.DatabaseOperation(dbId, "reload")
	if err != nil {
		return err
	}

	if err = c.attachToOperationLogs(op); err != nil {
		return err
	}

	return nil
}

// rename

// replicate

// restart

// tunnel

// url

// versions
