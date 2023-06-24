package internal

import (
	"log"

	"github.com/urfave/cli/v2"
)

type Resource struct {
	ID int64
}

func (c *Config) Logs(ctx *cli.Context) {
	appId := ctx.Value("app").(int64)
	dbId := ctx.Value("database").(int64)
	environmentId := ctx.Value("environment").(int64)

	if appId == 0 && dbId == 0 {
		log.Fatal("error - neither app nor database ids were provided, cannot continue")
	}

	if environmentId > 0 {
		// doesn't get used unless flag is passed in
		_, err := c.client.GetEnvironment(environmentId)
		if err != nil {
			log.Fatal(err)
		}
	}

	r := Resource{}
	if appId > 0 {
		app, err := c.client.GetApp(appId)
		if err != nil {
			log.Fatal(err)
		}
		if environmentId > 0 && app.EnvironmentID != environmentId {
			log.Fatal("error - app's environment and environment param do not match")
		}
		r.ID = app.ID
	}

	if dbId > 0 {
		db, err := c.client.GetDatabase(dbId)
		if err != nil {
			log.Fatal(err)
		}
		if environmentId > 0 && db.EnvironmentID != environmentId {
			log.Fatal("error - app's environment and environment param do not match")
		}
		r.ID = db.ID
	}
}
