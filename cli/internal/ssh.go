package internal

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func (c *Config) generatePublicPrivateKey() ([]byte, []byte, error) {
	tmp, err := os.MkdirTemp("", "sshkeygen-data")
	if err != nil {
		return nil, nil, err
	}
	defer os.RemoveAll(tmp)

	cmd := exec.Command(c.sshKeygenPath, "-t", "rsa", "-N", "", "-f", fmt.Sprintf("%s/id_rsa", tmp))
	_, err = cmd.Output()
	if err != nil {
		return nil, nil, err
	}

	privateKey, err := os.ReadFile(fmt.Sprintf("%s/id_rsa", tmp))
	if err != nil {
		return nil, nil, err
	}

	publicKey, err := os.ReadFile(fmt.Sprintf("%s/id_rsa.pub", tmp))
	if err != nil {
		return nil, nil, err
	}

	return publicKey, privateKey, nil
}

func (c *Config) AptibleSSH(publicKey []byte, privateKey []byte, certString, host, hostKey, user, token string, port int64) error {
	var err error

	if err = CheckHostPortAccessible(host, fmt.Sprintf("%d", port)); err != nil {
		return err
	}

	tmp, err := os.MkdirTemp("", "ssh-session-data")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	knownHostsPath := fmt.Sprintf("%s/known_hosts", tmp)
	if err = os.WriteFile(fmt.Sprintf("%s/id_rsa", tmp), privateKey, 0600); err != nil {
		return err
	}
	if err = os.WriteFile(fmt.Sprintf("%s/id_rsa.pub", tmp), publicKey, 0600); err != nil {
		return err
	}
	if err = os.WriteFile(fmt.Sprintf("%s/id_rsa-cert.pub", tmp), []byte(certString), 0600); err != nil {
		return err
	}
	if err = os.WriteFile(knownHostsPath, []byte(fmt.Sprintf("[%s]:%d %s", host, port, hostKey)), 0600); err != nil {
		return err
	}
	if err = os.Setenv("ACCESS_TOKEN", token); err != nil {
		return err
	}

	cmd := exec.Command(
		c.sshPath, fmt.Sprintf("%s@%s", user, host),
		"-i", fmt.Sprintf("%s/id_rsa", tmp),
		"-p", fmt.Sprintf("%d", port),
		"-o", "TCPKeepAlive=yes",
		"-o", "KeepAlive=yes",
		"-o", "ServerAliveInterval=60",
		"-o", "ControlMaster=no",
		"-o", "ControlPath=none",
		"-o", "SendEnv=ACCESS_TOKEN",
		"-o", "IdentitiesOnly=yes",
		"-o", fmt.Sprintf("UserKnownHostsFile=%s", knownHostsPath),
		"-o", "StrictHostKeyChecking=yes",
		"-T",
	)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = io.MultiWriter(os.Stdout, stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr, stderr)
	if err = cmd.Run(); err != nil {
		return err
	}

	return nil
}
