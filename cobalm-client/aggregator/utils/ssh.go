package utils

import (
	"aggregator/config"
	"aggregator/model"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

func GetDataFromRemote(cfg config.Config, t model.Target, outputPath string) error {
	// Create ssh config
	clientConfig := &ssh.ClientConfig{
		User: t.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(t.Pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to host
	client, err := ssh.Dial("tcp", t.NetIp, clientConfig)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// Prepare output file
	f, err := CreateOrReplaceFile(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Run ssh command
	session.Stdout = f
	session.Stdin = os.Stdin
	session.Stderr = os.Stderr

	return session.Run(fmt.Sprintf("cat %s; rm -f %s", t.File, t.File))
}
