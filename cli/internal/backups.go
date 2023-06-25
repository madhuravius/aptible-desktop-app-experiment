package internal

import (
	"fmt"
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

func (c *Config) ListBackups(ctx *cli.Context) error {
	var envs []aptible.Environment
	var err error

	environmentId := ctx.Value("environment").(int64)
	dbId := ctx.Value("database").(int64)

	var dbs []aptible.Database
	if dbId != 0 {
		db, err := c.client.GetDatabase(dbId)
		if err != nil {
			return err
		}
		dbs = []aptible.Database{db}
	} else {
		if environmentId != 0 {
			environment, err := c.client.GetEnvironment(environmentId)
			if err != nil {
				return err
			}
			envs = []aptible.Environment{environment}
		} else {
			envs, err = c.client.GetEnvironments()
			if err != nil {
				return err
			}
		}
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
