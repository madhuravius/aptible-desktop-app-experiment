package internal

import (
	"fmt"
	"log"

	"github.com/aptible/go-deploy/aptible"
	"github.com/urfave/cli/v2"
)

func ListLogDrains(cCtx *cli.Context) {
	var err error
	client, err := Client(cCtx)
	if err != nil {
		log.Fatal(err)
	}

	environmentId := cCtx.Value("environment").(int64)
	var envs []aptible.Environment
	if environmentId > 0 {
		environment, err := client.GetEnvironment(environmentId)
		if err != nil {
			log.Fatal(err)
		}
		envs = []aptible.Environment{environment}
	} else {
		envs, err = client.GetEnvironments()
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, env := range envs {
		dbs, err := client.GetDatabases(env.ID)
		if err != nil {
			log.Fatal(err)
		}
		if len(dbs) == 0 {
			continue
		}

		fmt.Printf("=== %s\n", env.Handle)
		for _, db := range dbs {
			fmt.Println(db.Handle)
		}
	}
}
