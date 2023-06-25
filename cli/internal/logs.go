package internal

import (
	"errors"
	"github.com/aptible/go-deploy/aptible"

	"github.com/urfave/cli/v2"
)

func (c *Config) Logs(ctx *cli.Context) error {
	var err error
	appId := ctx.Value("app").(int64)
	dbId := ctx.Value("database").(int64)
	environmentId := ctx.Value("environment").(int64)

	if appId == 0 && dbId == 0 {
		return errors.New("error - neither app nor database ids were provided, cannot continue")
	}

	if environmentId != 0 {
		// doesn't get used unless flag is passed in
		_, err = c.client.GetEnvironment(environmentId)
		if err != nil {
			return err
		}
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
