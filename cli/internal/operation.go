package internal

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/urfave/cli/v2"
)

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

	publicKey, privateKey, err := generatePublicPrivateKey()
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
	sshPortalOp, err := c.client.CreateSSHPortalConnectionOperation(op.EnvironmentID, opId, string(publicKey))
	if err != nil {
		return err
	}

	fmt.Print(Green(fmt.Sprintf("Streaming logs for running %s #%d on %s...\n", op.Type, op.ID, op.Handle)))
	err = AptibleSSH(privateKey, sshPortalOp.Certificate, stack.PortalHost, stack.HostKey, sshPortalOp.SSHUser, c.token, stack.PortalPort)
	if err != nil {
		return err
	}

	return nil
}
