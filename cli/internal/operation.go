package internal

import (
	"errors"
	"fmt"
	"github.com/aptible/go-deploy/aptible"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/urfave/cli/v2"
)

func (c *Config) attachToOperationLogs(op aptible.Operation) error {
	publicKey, privateKey, err := c.generatePublicPrivateKey()
	if err != nil {
		return err
	}

	environment, err := c.client.GetEnvironment(op.EnvironmentID)
	if err != nil {
		return err
	}

	stack, err := c.client.GetStack(environment.StackID)
	if err != nil {
		return err
	}

	// create SSH Portal connection operation
	sshPortalOp, err := c.client.CreateSSHPortalConnectionOperation(op.EnvironmentID, op.ID, string(publicKey))
	if err != nil {
		return err
	}

	fmt.Print(Green(fmt.Sprintf("Streaming logs for running %s #%d on %s...\n", op.Type, op.ID, op.Handle)))
	err = c.AptibleSSH(publicKey, privateKey, sshPortalOp.Certificate, stack.PortalHost, stack.HostKey, sshPortalOp.SSHUser, c.token, stack.PortalPort)
	if err != nil {
		return err
	}

	return nil
}

func GenOperationsCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "operation:follow",
			Usage: "This command follows the logs of a running Operation.",
			Action: func(ctx *cli.Context) error {
				c := NewConfigF(ctx)
				if err := c.OperationFollow(ctx); err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
		{
			Name:  "operation:logs",
			Usage: "This command displays logs for a given Operation.",
			Action: func(ctx *cli.Context) error {
				c := NewConfigF(ctx)
				if err := c.OperationFollow(ctx); err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
	}
}

func (c *Config) OperationFollow(ctx *cli.Context) error {
	// setup SSH and ensure access
	if len(ctx.Args().Slice()) == 0 {
		return errors.New("missing operation id argument to continue")
	}

	opIdRaw := ctx.Args().Get(0)
	opId, err := strconv.ParseInt(opIdRaw, 10, 64)
	if err != nil {
		return err
	}
	if opId == 0 {
		return errors.New("missing operation id argument to continue")
	}

	// Streaming logs for running configure #4644 on test-123...
	op, err := c.client.GetOperation(opId)
	if err != nil {
		return err
	}

	if op.Status == "failed" || op.Status == "succeeded" {
		fmt.Print("This operation has already succeeded. ")
		fmt.Println("Run the following command to retrieve the operation's logs:")
		fmt.Printf("aptible operation:logs %d\n", op.ID)
		return nil
	}

	if err = c.attachToOperationLogs(op); err != nil {
		return err
	}

	return nil
}

func (c *Config) OperationLogs(ctx *cli.Context) error {
	if len(ctx.Args().Slice()) == 0 {
		return errors.New("missing operation id argument to continue")
	}

	opIdRaw := ctx.Args().Get(0)
	opId, err := strconv.ParseInt(opIdRaw, 10, 64)
	if err != nil {
		return err
	}
	if opId == 0 {
		return errors.New("missing operation id argument to continue")
	}

	op, err := c.client.GetOperation(opId)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/operations/%d/logs", c.apiHost, op.ID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	reqLogs, err := http.NewRequest("GET", string(body), nil)
	if err != nil {
		return err
	}

	respLogs, err := client.Do(reqLogs)
	if err != nil {
		return err
	}

	logsBody, err := io.ReadAll(respLogs.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(logsBody))

	return nil
}
