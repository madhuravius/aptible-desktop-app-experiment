package internal

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
)

func GenBackupsCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name: "backup:list",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "environment",
					Usage: "Specify an environment to run your backup:list command on",
				},
				&cli.StringFlag{
					Name:  "database",
					Usage: "Specify an database to run your backup:list command on",
				},
			},
			Usage: "This command lists all Database Backups for a given Database.",
			Action: func(ctx *cli.Context) error {
				c := NewConfigF(ctx)
				if err := c.ListBackups(ctx); err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
	}
}

func (c *Config) ListBackups(ctx *cli.Context) error {
	var err error

	dbs, err := c.getDatabasesFromFlags(ctx)
	if err != nil {
		return err
	}

	for _, db := range dbs {
		backups, err := c.client.GetBackups(db.ID)
		if err != nil {
			return err
		}

		for _, backup := range backups {
			date, err := translateDateToStdOut(backup.CreatedAt)
			if err != nil {
				return err
			}
			fmt.Printf("%d: %s, %s, ", backup.ID, date, backup.Region)
			if backup.Manual {
				fmt.Print("manual")
			} else {
				fmt.Print("automatic")
			}

			if backup.Copy != nil {
				fmt.Print(", copy")
			}
		}
	}

	return nil
}
