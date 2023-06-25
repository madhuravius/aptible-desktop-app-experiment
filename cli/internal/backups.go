package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/aptible/go-deploy/aptible"
	"github.com/urfave/cli/v2"
)

func translateDateToBackupStdOut(date string) (string, error) {
	// 2023-06-23T05:55:20.000Z - in
	// 2023-05-25 05:01:07 UTC - out
	parsedDate, err := time.Parse("2006-01-02T15:04:05.000Z", date)
	if err != nil {
		return "", err
	}

	return parsedDate.Format("2006-01-02 15:04:05 UTC"), nil
}

func GenBackupsCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name: "backup:list",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "environment",
					Usage: "Specify an environment to run your backup:list command on",
				},
				&cli.Int64Flag{
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
	var envs []aptible.Environment
	var err error

	envs, err = c.getEnvironmentsFromFlags(ctx)
	if err != nil {
		return err
	}

	dbId := ctx.Value("database").(int64)

	var dbs []aptible.Database
	if dbId != 0 {
		db, err := c.client.GetDatabase(dbId)
		if err != nil {
			return err
		}
		dbs = []aptible.Database{db}
	} else {
		for _, env := range envs {
			dbsInEnv, err := c.client.GetDatabases(env.ID)
			if err != nil {
				return err
			}
			if len(dbsInEnv) == 0 {
				continue
			}
			dbs = append(dbs, dbsInEnv...)
		}
	}

	for _, db := range dbs {
		backups, err := c.client.GetBackups(db.ID)
		if err != nil {
			return err
		}

		for _, backup := range backups {
			date, err := translateDateToBackupStdOut(backup.CreatedAt)
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
