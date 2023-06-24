package internal

import (
	"os"

	"github.com/aptible/go-deploy/aptible"
	"github.com/urfave/cli/v2"
)

type Config struct {
	client *aptible.Client
}

func NewConfig(ctx *cli.Context) (*Config, error) {
	client, err := Client(ctx)
	if err != nil {
		return nil, err
	}
	return &Config{
		client: client,
	}, nil
}

func Client(ctx *cli.Context) (*aptible.Client, error) {
	token := ctx.Value("token").(string)
	apiHost := ctx.Value("api-host").(string)

	os.Setenv("APTIBLE_ACCESS_TOKEN", token)
	os.Setenv("APTIBLE_API_ROOT_URL", apiHost)

	client, err := aptible.SetUpClient()
	if err != nil {
		return nil, err
	}
	return client, err
}

func (c *Config) getEnvironmentsFromFlags(ctx *cli.Context) ([]aptible.Environment, error) {
	var err error
	environmentId := ctx.Value("environment").(int64)
	var envs []aptible.Environment
	if environmentId > 0 {
		environment, err := c.client.GetEnvironment(environmentId)
		if err != nil {
			return nil, err
		}
		envs = []aptible.Environment{environment}
	} else {
		envs, err = c.client.GetEnvironments()
		if err != nil {
			return nil, err
		}
	}
	return envs, nil
}
