package internal

import (
	"fmt"
	"os"

	"github.com/aptible/go-deploy/aptible"
	"github.com/urfave/cli/v2"
)

var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
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

	// todo - find a way to bypass this, this is pretty bad
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
	if environmentId != 0 {
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

// Color - adjust color based on supplied codes, supplied from:
// https://gist.github.com/ik5/d8ecde700972d4378d87?permalink_comment_id=3074524#gistcomment-3074524
func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}
