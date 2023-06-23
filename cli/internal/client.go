package internal

import (
	"os"

	"github.com/aptible/go-deploy/aptible"
	"github.com/urfave/cli/v2"
)

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
