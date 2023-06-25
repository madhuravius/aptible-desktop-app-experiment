package internal

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

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
	client        *aptible.Client
	token         string
	apiHost       string
	sshPath       string
	sshKeygenPath string
}

func NewConfigF(ctx *cli.Context) *Config {
	c, err := NewConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func NewConfig(ctx *cli.Context) (*Config, error) {
	token := ctx.Value("token").(string)
	apiHost := ctx.Value("api-host").(string)

	sshPath := os.Getenv("SSH_PATH")
	sshKeygenPath := os.Getenv("SSH_KEYGEN_PATH")

	client, err := Client(token, apiHost)
	if err != nil {
		return nil, err
	}
	return &Config{
		client:        client,
		token:         token,
		apiHost:       apiHost,
		sshPath:       sshPath,
		sshKeygenPath: sshKeygenPath,
	}, nil
}

func Client(token, apiHost string) (*aptible.Client, error) {
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

	// if there is an error, we will skip it as we defer to whatever list is provided instead
	environmentId, _ := c.getEnvironmentIDFromFlags(ctx)

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

func (c *Config) getEnvironmentIDFromFlags(ctx *cli.Context) (int64, error) {
	rawEnvIdOrHandle := ctx.Value("environment").(string)
	if envId, err := strconv.ParseInt(rawEnvIdOrHandle, 10, 64); err != nil {
		envs, envsErr := c.client.GetEnvironments()
		if envsErr != nil {
			return 0, fmt.Errorf("could not query environments to get environmentId: %s", err.Error())
		}
		for _, env := range envs {
			if env.Handle == rawEnvIdOrHandle {
				return env.ID, nil
			}
		}
	} else {
		return envId, nil
	}
	return 0, fmt.Errorf("specified environment does not exist: %s", rawEnvIdOrHandle)
}

func (c *Config) getAppIDFromFlags(ctx *cli.Context) (int64, error) {
	environments, err := c.getEnvironmentsFromFlags(ctx)
	if err != nil {
		return 0, err
	}

	rawAppIdOrHandle := ctx.Value("app").(string)
	if appId, err := strconv.ParseInt(rawAppIdOrHandle, 10, 64); err != nil {
		for _, env := range environments {
			apps, appsErr := c.client.GetApps(env.ID)
			if appsErr != nil {
				return 0, fmt.Errorf("could not query apps to get appId from handle: %s", err.Error())
			}
			for _, app := range apps {
				if app.Handle == rawAppIdOrHandle {
					return app.ID, nil
				}
			}
		}
	} else {
		return appId, nil
	}
	return 0, fmt.Errorf("specified app does not exist: %s", rawAppIdOrHandle)
}

func (c *Config) getDatabaseIDFromFlags(ctx *cli.Context) (int64, error) {
	environments, err := c.getEnvironmentsFromFlags(ctx)
	if err != nil {
		return 0, err
	}

	rawDbIdOrHandle := ctx.Value("database").(string)
	if dbId, err := strconv.ParseInt(rawDbIdOrHandle, 10, 64); err != nil {
		for _, env := range environments {
			dbs, dbsErr := c.client.GetDatabases(env.ID)
			if dbsErr != nil {
				return 0, fmt.Errorf("could not query databases to get databaseId from handle: %s", err.Error())
			}
			for _, db := range dbs {
				if db.Handle == rawDbIdOrHandle {
					return db.ID, nil
				}
			}
		}
	} else {
		return dbId, nil
	}
	return 0, fmt.Errorf("specified database does not exist: %s", rawDbIdOrHandle)
}

func (c *Config) getDatabasesFromFlags(ctx *cli.Context) ([]aptible.Database, error) {
	rawDbIdOrHandle := ctx.Value("database").(string)

	// if there is an error, we will skip it as we defer to whatever list is provided instead
	var dbs []aptible.Database
	if dbId, err := strconv.ParseInt(rawDbIdOrHandle, 10, 64); err == nil {
		db, err := c.client.GetDatabase(dbId)
		if err != nil {
			return nil, err
		}
		dbs = []aptible.Database{db}
	} else {
		envs, err := c.getEnvironmentsFromFlags(ctx)
		if err != nil {
			return nil, err
		}

		for _, env := range envs {
			dbResults, dbsErr := c.client.GetDatabases(env.ID)
			if dbsErr != nil {
				return nil, fmt.Errorf("could not query databases to collect for environment: %s", err.Error())
			}
			dbs = append(dbs, dbResults...)
		}
	}

	return dbs, nil
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

func CheckHostPortAccessible(host, port string) error {
	// check if host is available / port open
	checkConn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 1*time.Second)
	if err != nil {
		return err
	}
	if checkConn != nil {
		_ = checkConn.Close()
	} else {
		return errors.New("error - unable to connect to remote host")
	}

	return nil
}

func translateDateToStdOut(date string) (string, error) {
	// 2023-06-23T05:55:20.000Z - in
	// 2023-05-25 05:01:07 UTC - out
	parsedDate, err := time.Parse("2006-01-02T15:04:05.000Z", date)
	if err != nil {
		return "", err
	}

	return parsedDate.Format("2006-01-02 15:04:05 UTC"), nil
}
