package internal

import (
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

func SetupSSHClient(user, password, host, publicKeyPath, privateKeyPath string) {
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal(err)
	}
	signer, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		log.Fatal(err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	conn, err := ssh.Dial("tcp", host, config)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var session *ssh.Session
	session, err = conn.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
}
