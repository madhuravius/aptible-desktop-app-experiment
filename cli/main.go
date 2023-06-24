package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aptible/aptible-desktop-app/cli/internal"
	"github.com/urfave/cli/v2" // imports as package "cli"
)

const logo = `
      ..'...''..             .','.                          ... .',.         .,'.
    .;oxo;.'cddc'.          .:0NXo.                  .od,  .:kc.'kK:        .lXx.
  .;oxo:'....,lxxl,.        ,OXk0Kc.   ......'''..  .:KNo....,. 'OXc.'''..  .oNk.   ..''...
 'oxo;.,lc.,l:',cxx:.      .xNx'cX0,   'x0xxxxkOOd,.:ONW0d,.l0l.'0W0xxxkOx,..oNk. .cxkxxkOo'.
 ,l;.,lxOl.;xkd:.'cc.     .lXO, .dNx.  ,0W0:....lKK:.:KNo...dWx.'0WO;..'oX0;.oNk..xXx,..'oKO,
 ..,lxkkOl.;xOkxd:...    .:KWOlclxXNo. ,0No.    .xWd.,0Xc. .dWx.'0Nl.   .kNo.oNk.:KNOddddx00c.
 .lxxc;oOl.;xkc;lxd;.    ,ON0xxxxxOXXc.,0Wx.    'ONo.,0Xc. .dWx.'OWd.   ,0Nl.oNk.,0Xo'''';lc'
 ,oc'..lOl.;xk:..,ll.   .xNO'     .dN0,,0WXxlccoOKx' 'kNOl'.dNx.'OWKxlco0Xd..oNk..:OOdcclk0o.
 ...  .,:'..;;.   ...   .cl,.      .cl,;0Xdcodxdl,.  .'col'.,l,..:l::odol,. .,l;. ..;lddoc,.
                                       ,OK:.
                                       .,;.
`

const desc = `aptible is a command line interface to the Aptible.com platform.

It allows users to manage authentication, application launch, deployment, logging, and more.
To read more, use the docs command to view Aptible's help on the web.`

func genConfig(ctx *cli.Context) *internal.Config {
	c, err := internal.NewConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func main() {
	app := &cli.App{
		Name:  "aptible",
		Usage: "aptible cli",
		Commands: []*cli.Command{
			{
				Name:  "about",
				Usage: "print some information about the CLI aptible CLI",
				Action: func(_ *cli.Context) error {
					fmt.Printf("%s\n%s\n", logo, desc)
					return nil
				},
			},
			{
				Name: "apps",
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:  "environment",
						Value: 0,
						Usage: "Specify an environment to run your apps:list command on",
					},
				},
				Usage: "This command lists Apps in an Environment.",
				Action: func(ctx *cli.Context) error {
					c := genConfig(ctx)
					if err := c.ListApps(ctx); err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name: "backup:list",
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:  "environment",
						Value: 0,
						Usage: "Specify an environment to run your backup:list command on",
					},
					&cli.Int64Flag{
						Name:  "database",
						Value: 0,
						Usage: "Specify an database to run your backup:list command on",
					},
				},
				Usage: "This command lists all Database Backups for a given Database.",
				Action: func(ctx *cli.Context) error {
					c := genConfig(ctx)
					if err := c.ListBackups(ctx); err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name: "db:list",
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:  "environment",
						Value: 0,
						Usage: "Specify an environment to run your db:list command on",
					},
				},
				Usage: "This command lists Databases in an Environment.",
				Action: func(ctx *cli.Context) error {
					c := genConfig(ctx)
					if err := c.ListDatabases(ctx); err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
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
					c := genConfig(ctx)
					if err := c.ListEndpoints(ctx); err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name: "environment:list",
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:  "environment",
						Value: 0,
						Usage: "Specify an environment to run your environment:list command on",
					},
				},
				Usage: "This command lists all Environments.",
				Action: func(ctx *cli.Context) error {
					c := genConfig(ctx)
					if err := c.ListEnvironments(ctx); err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name: "log_drain:list",
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:  "environment",
						Value: 0,
						Usage: "Specify an environment to run your log_drain:list command on",
					},
				},
				Usage: "This command lists all Log Drains.",
				Action: func(ctx *cli.Context) error {
					c := genConfig(ctx)
					if err := c.ListLogDrains(ctx); err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name: "metric_drain:list",
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:  "environment",
						Value: 0,
						Usage: "Specify an environment to run your metric_drain:list command on",
					},
				},
				Usage: "This command lists all Metric Drains.",
				Action: func(ctx *cli.Context) error {
					c := genConfig(ctx)
					if err := c.ListMetricDrains(ctx); err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name: "logs",
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:  "environment",
						Value: 0,
						Usage: "Specify an environment to run your logs command on",
					},
					&cli.Int64Flag{
						Name:  "app",
						Value: 0,
						Usage: "Specify an app to run your logs command on",
					},
					&cli.Int64Flag{
						Name:  "database",
						Value: 0,
						Usage: "Specify an database to run your logs command on",
					},
				},
				Usage: "This command lets you access real-time logs for an App or Database.",
				Action: func(ctx *cli.Context) error {
					c := genConfig(ctx)
					if err := c.Logs(ctx); err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "token",
				Usage: "Specify a api token to use when running requests",
			},
			&cli.StringFlag{
				Name:  "api-host",
				Usage: "Specify a api-host you want your commands to run against",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
