package utils

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func ConnectSsh(host string, user string, pass string, folder string) error {
	clientConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host, clientConfig)
	if err != nil {
		return err
	}

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	err = session.RequestPty("linux", 80, 40, modes)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	session.Stdout = &buffer

	fmt.Printf("Running commands\n")
	err = session.Run(fmt.Sprintf("cat %s; rm %s", folder, folder))

	if err != nil {
		fmt.Printf("Got an error while executing: %s", err.Error())
		return err
	}

	fmt.Printf("cat function returned:\n%s\nprint ended", buffer.String())

	return nil
}
