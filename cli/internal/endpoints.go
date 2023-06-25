package internal

import (
	"fmt"
	"log"

	"github.com/aptible/go-deploy/aptible"
	"github.com/urfave/cli/v2"
)

func translateProtocolToReadable(protocol string) string {
	if protocol == "http_proxy_protocol" || protocol == "http" {
		return "https"
	}
	return protocol
}

func translateIPWhitelistToReadable(whitelist []string) string {
	if len(whitelist) == 0 {
		return "all traffic"
	}
	var output string
	for idx, ip := range whitelist {
		output += ip
		if idx != len(whitelist)-1 && idx > 0 {
			output += ","
		}
	}
	return output
}

func translateContainerPortToReadable(port int64) string {
	if port == 0 {
		return "default"
	}
	return fmt.Sprintf("%v", port)
}

func printEndpoint(endpoint aptible.Endpoint) {
	fmt.Printf("Id: %d\n", endpoint.ID)
	fmt.Printf("Hostname: %s\n", endpoint.ExternalHost)
	fmt.Printf("Status: %s\n", endpoint.Status)
	fmt.Printf("Created At: %s\n", endpoint.CreatedAt)
	fmt.Printf("Type: %s\n", translateProtocolToReadable(endpoint.Type))
	fmt.Printf("Port: %s\n", translateContainerPortToReadable(endpoint.ContainerPort))
	fmt.Printf("Internal: %v\n", endpoint.Internal)
	fmt.Printf("IP Whitelist: %s\n", translateIPWhitelistToReadable(endpoint.IPWhitelist))
	fmt.Printf("Default Domain Enabled: %v\n", endpoint.Default)
	if endpoint.AcmeChallenges == nil || len(endpoint.AcmeChallenges) == 0 {
		fmt.Println("Managed TLS Enabled: false")
	} else {
		fmt.Println("Managed TLS Enabled: true")
		for _, acmeChallenge := range endpoint.AcmeChallenges {
			fmt.Printf("Managed TLS Domain: %s\n", acmeChallenge.To)
			fmt.Printf("Managed TLS DNS Challenge Hostname: %s\n", acmeChallenge.From)
			fmt.Printf("Managed TLS Status: %s\n", acmeChallenge.Status)
		}
	}
	fmt.Printf("Service: %s\n\n", endpoint.Service.ProcessType)
}

func printEndpointsIfInServiceEndpoints(endpoint aptible.Endpoint, services []aptible.Service) {
	for _, service := range services {
		if service.ResourceID != endpoint.Service.ResourceID {
			continue
		}
		printEndpoint(endpoint)
	}
}

func GenEndpointsCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name: "endpoints:list",
			Flags: []cli.Flag{
				&cli.Int64Flag{
					Name:  "app",
					Value: 0,
					Usage: "Specify an app to run your endpoints:list command on",
				},
				&cli.Int64Flag{
					Name:  "database",
					Value: 0,
					Usage: "Specify a database to run your endpoints:list command on",
				},
				&cli.Int64Flag{
					Name:  "environment",
					Value: 0,
					Usage: "Specify an environment to run your endpoints:list command on",
				},
			},
			Usage: "This command lists all Endpoints.",
			Action: func(ctx *cli.Context) error {
				c := NewConfigF(ctx)
				if err := c.ListEndpoints(ctx); err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
	}
}

func (c *Config) ListEndpoints(ctx *cli.Context) error {
	var err error
	var app aptible.App
	var db aptible.Database

	envs, err := c.getEnvironmentsFromFlags(ctx)
	if err != nil {
		return err
	}

	appId := ctx.Value("app").(int64)
	if appId != 0 {
		app, err = c.client.GetApp(appId)
		if err != nil {
			return err
		}
	}

	dbId := ctx.Value("database").(int64)
	if dbId != 0 {
		db, err = c.client.GetDatabase(dbId)
		if err != nil {
			return err
		}
	}

	for _, env := range envs {
		endpoints, err := c.client.GetEndpoints(env.ID)
		if err != nil {
			return err
		}
		if len(endpoints) == 0 {
			continue
		}

		if app.ID != 0 && app.EnvironmentID == env.ID ||
			db.ID != 0 && db.EnvironmentID == env.ID ||
			app.ID == 0 && db.ID == 0 {
			fmt.Printf("=== %s\n", env.Handle)
		}
		for _, endpoint := range endpoints {
			if app.ID != 0 {
				printEndpointsIfInServiceEndpoints(endpoint, app.Services)
			} else if db.ID != 0 {
				printEndpointsIfInServiceEndpoints(endpoint, []aptible.Service{db.Service})
			} else {
				printEndpoint(endpoint)
			}
		}
	}
	return nil
}
