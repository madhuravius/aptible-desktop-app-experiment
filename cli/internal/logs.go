package internal

import (
	"errors"
	"github.com/aptible/go-deploy/aptible"
	"log"

	"github.com/urfave/cli/v2"
)

func GenLogsCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name: "logs",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "environment",
					Usage: "Specify an environment to run your logs command on",
				},
				&cli.Int64Flag{
					Name:  "app",
					Usage: "Specify an app to run your logs command on",
				},
				&cli.StringFlag{
					Name:  "database",
					Usage: "Specify an database to run your logs command on",
				},
			},
			Usage: "This command lets you access real-time logs for an App or Database.",
			Action: func(ctx *cli.Context) error {
				c := NewConfigF(ctx)
				if err := c.Logs(ctx); err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
	}
}

func (c *Config) Logs(ctx *cli.Context) error {
	var err error
	appId := ctx.Value("app").(int64)

	dbId, _ := c.getDatabaseIDFromFlags(ctx)
	environmentId, _ := c.getEnvironmentIDFromFlags(ctx)

	if appId == 0 && dbId == 0 {
		return errors.New("error - neither app nor database ids were provided, cannot continue")
	}

	var op aptible.Operation
	if appId != 0 {
		app, err := c.client.GetApp(appId)
		if err != nil {
			return err
		}
		if environmentId != 0 && app.EnvironmentID != environmentId {
			return errors.New("error - app's environment and environment param do not match")
		}
		op, err = c.client.CreateAppLogsOperation(appId)
		if err != nil {
			return err
		}
	}

	if dbId != 0 {
		db, err := c.client.GetDatabase(dbId)
		if err != nil {
			return err
		}
		if environmentId != 0 && db.EnvironmentID != environmentId {
			return errors.New("error - app's environment and environment param do not match")
		}
		op, err = c.client.CreateDatabaseLogsOperation(dbId)
		if err != nil {
			return err
		}
	}

	if err = c.attachToOperationLogs(op); err != nil {
		return err
	}

	return nil
}
