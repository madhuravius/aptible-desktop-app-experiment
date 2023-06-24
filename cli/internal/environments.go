package internal

import (
	"fmt"
	"log"

	"github.com/aptible/go-deploy/aptible"
	"github.com/urfave/cli/v2"
)

func ListEnvironments(cCtx *cli.Context) {
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
		fmt.Printf("Got environment successfully - %s (%d)\n", environment.Handle, environment.ID)
	} else {
		envs, err = client.GetEnvironments()
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, env := range envs {
		fmt.Println(env.Handle)
	}
}
