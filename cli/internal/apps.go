package internal

import (
	"fmt"
	"github.com/aptible/go-deploy/aptible"
	"github.com/urfave/cli/v2"
	"log"
)

func GenAppCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name: "apps",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "environment",
					Usage: "Specify an environment to run your apps command on",
				},
			},
			Usage: "This command lists Apps in an Environment.",
			Action: func(ctx *cli.Context) error {
				c := NewConfigF(ctx)
				if err := c.ListApps(ctx); err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
		{
			Name: "rebuild",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "environment",
					Usage: "Specify an environment to run your app:rebuild command on",
				},
				&cli.StringFlag{
					Name:     "app",
					Required: true,
					Usage:    "Specify an app to run your apps:rebuild command on",
				},
			},
			Usage: "This command rebuilds an App and restarts its Services.",
			Action: func(ctx *cli.Context) error {
				c := NewConfigF(ctx)
				if err := c.RebuildApp(ctx); err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
		{
			Name: "restart",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "environment",
					Usage: "Specify an environment to run your app:rebuild command on",
				},
				&cli.StringFlag{
					Name:     "app",
					Required: true,
					Usage:    "Specify an app to run your apps:rebuild command on",
				},
			},
			Usage: "This command restarts an App and all its associated Services.",
			Action: func(ctx *cli.Context) error {
				c := NewConfigF(ctx)
				if err := c.RestartApp(ctx); err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
		{
			Name: "services",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "environment",
					Usage: "Specify an environment to run your app:rebuild command on",
				},
				&cli.StringFlag{
					Name:     "app",
					Required: true,
					Usage:    "Specify an app to run your apps:rebuild command on",
				},
			},
			Usage: "This command lists all Services for a given App.",
			Action: func(ctx *cli.Context) error {
				c := NewConfigF(ctx)
				if err := c.GetServices(ctx); err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
	}
}

func (c *Config) ListApps(ctx *cli.Context) error {
	var envs []aptible.Environment
	var err error

	envs, err = c.getEnvironmentsFromFlags(ctx)
	if err != nil {
		return err
	}

	for _, env := range envs {
		apps, err := c.client.GetApps(env.ID)
		if err != nil {
			return err
		}
		if len(apps) == 0 {
			continue
		}

		fmt.Printf("=== %s\n", env.Handle)
		for _, app := range apps {
			fmt.Println(app.Handle)
		}
	}
	return nil
}

// config

// deploy
func (c *Config) RebuildApp(ctx *cli.Context) error {
	appId, err := c.getAppIDFromFlags(ctx)
	if err != nil {
		return err
	}

	op, err := c.client.AppOperation(appId, "rebuild")
	if err != nil {
		return err
	}

	if err = c.attachToOperationLogs(op); err != nil {
		return err
	}

	return nil
}

func (c *Config) RestartApp(ctx *cli.Context) error {
	appId, err := c.getAppIDFromFlags(ctx)
	if err != nil {
		return err
	}

	op, err := c.client.AppOperation(appId, "restart")
	if err != nil {
		return err
	}

	if err = c.attachToOperationLogs(op); err != nil {
		return err
	}

	return nil
}

func (c *Config) GetServices(ctx *cli.Context) error {
	appId, err := c.getAppIDFromFlags(ctx)
	if err != nil {
		return err
	}

	app, err := c.client.GetApp(appId)
	if err != nil {
		return err
	}

	for _, service := range app.Services {
		createdAt, err := translateDateToStdOut(service.CreatedAt)
		if err != nil {
			return err
		}
		fmt.Printf("Id: %d\n", service.ID)
		fmt.Printf("Service: %s\n", service.ProcessType)
		fmt.Printf("Created at: %s\n", createdAt)
		fmt.Printf("Command: %s\n", service.Command)
		fmt.Printf("Container Count: %d\n", service.ContainerCount)
		fmt.Printf("Container Size: %d\n", service.ContainerMemoryLimitMb)
		fmt.Printf("App: %s\n", app.Handle)
	}

	return nil
}

// services
