package internal

import (
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"

	"golang.org/x/crypto/ssh"
)

// keyString - from a given public key will check generated hosts value and verify: https://stackoverflow.com/a/63308243
// child stanza just to generate the match
func keyString(k ssh.PublicKey) string {
	return k.Type() + " " + base64.StdEncoding.EncodeToString(k.Marshal()) // e.g. "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTY...."
}

// trustedHostKeyCallback - see above comment form keyStanza for attribution and how this is meant to actually work
func trustedHostKeyCallback(trustedKey string) ssh.HostKeyCallback {
	return func(_ string, _ net.Addr, k ssh.PublicKey) error {
		ks := keyString(k)
		if trustedKey != ks {
			return fmt.Errorf("SSH-key verification: expected %q but got %q", trustedKey, ks)
		}
		return nil
	}
}

func generatePublicPrivateKey() ([]byte, []byte, error) {
	tmp, err := os.MkdirTemp("", "sshkeygen-data")
	if err != nil {
		return nil, nil, err
	}
	defer os.RemoveAll(tmp)

	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-N", "", "-f", fmt.Sprintf("%s/id_rsa", tmp))
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

func AptibleSSH(privateKey []byte, certString, host, hostKey, user, token string, port int64) error {
	var err error

	if err = CheckHostPortAccessible(host, fmt.Sprintf("%d", port)); err != nil {
		return err
	}

	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return err
	}

	certBytes := []byte(certString)
	cert, _, _, _, err := ssh.ParseAuthorizedKey(certBytes)
	if err != nil {
		return err
	}

	certSigner, err := ssh.NewCertSigner(cert.(*ssh.Certificate), signer)
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User: user,
		// TODO - see: https://github.com/FiloSottile/yubikey-agent/blob/main/main.go#L263
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(certSigner)},
		HostKeyCallback: trustedHostKeyCallback(hostKey),
		Timeout:         30 * time.Second,
	}

	// actually ssh in
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		return err
	}
	defer conn.Close()

	var session *ssh.Session
	session, err = conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// set options for ssh connection on host
	if err = session.Setenv("ACCESS_TOKEN", token); err != nil {
		return err
	}

	return nil
}

/*
	1. check the resource we want is legit
	2. find the context for whatever thing we're trying to do:
		- operation logs (needs an existing operation to latch onto): https://github.com/aptible/aptible-cli/blob/2baa61e9ca55224d659d784fb4a8b14a3b7dbbb1/lib/aptible/cli/subcommands/operation.rb#L23
		- regular logs (logs operation): https://github.com/aptible/aptible-cli/blob/2baa61e9ca55224d659d784fb4a8b14a3b7dbbb1/lib/aptible/cli/subcommands/logs.rb#L26C15-L26C15
		- sshing to container (execute operation): https://github.com/aptible/aptible-cli/blob/2baa61e9ca55224d659d784fb4a8b14a3b7dbbb1/lib/aptible/cli/subcommands/ssh.rb#L52
*/
