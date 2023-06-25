package internal

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
)

func (c *Config) GetConfiguration(ctx *cli.Context) error {
	var err error

	appId := ctx.Value("app").(int64)
	environmentId := ctx.Value("environment").(int64)

	app, err := c.client.GetApp(appId)
	if err != nil {
		return err
	}

	if environmentId != 0 {
		environment, err := c.client.GetEnvironment(environmentId)
		if err != nil {
			return err
		}

		if app.EnvironmentID != environment.ID {
			return errors.New("error - app's environment and environment param do not match")
		}
	}

	if app.Env != nil {
		data, err := json.Marshal(app.Env)
		if err != nil {
			return err
		}
		out := make(map[string]interface{})
		json.Unmarshal(data, &out)

		for key, value := range out {
			fmt.Printf("%s=%s\n", key, value)
		}
	}

	return nil
}

// get
