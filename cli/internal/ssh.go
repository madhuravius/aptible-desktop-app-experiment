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

/*
	1. check the resource we want is legit
	2. find the context for whatever thing we're trying to do:
		- operation logs (needs an existing operation to latch onto): https://github.com/aptible/aptible-cli/blob/2baa61e9ca55224d659d784fb4a8b14a3b7dbbb1/lib/aptible/cli/subcommands/operation.rb#L23
		- regular logs (logs operation): https://github.com/aptible/aptible-cli/blob/2baa61e9ca55224d659d784fb4a8b14a3b7dbbb1/lib/aptible/cli/subcommands/logs.rb#L26C15-L26C15
		- sshing to container (execute operation): https://github.com/aptible/aptible-cli/blob/2baa61e9ca55224d659d784fb4a8b14a3b7dbbb1/lib/aptible/cli/subcommands/ssh.rb#L52
	3. check the default path of ~/home/.ssh or a preprovided path
	4. check if files are present where we expect them to exist:
		- ssh config file
		- private key file (`id_rsa`` by default)
		- public key file  (`id_rsa.pub` by default), or private key + .pub affix
		- we need the ACCESS_TOKEN, which is also sent here: https://github.com/aptible/aptible-cli/blob/master/lib/aptible/cli/helpers/operation.rb#L32
	5. once we have that, start the thing and pray
*/
