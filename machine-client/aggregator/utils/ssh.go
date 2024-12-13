package utils

import (
	"bytes"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

func GetFileFromRemote(host string, user string, pass string, filePath string) (bytes.Buffer, error) {
	clientConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host, clientConfig)
	if err != nil {
		return bytes.Buffer{}, err
	}

	session, err := client.NewSession()
	if err != nil {
		return bytes.Buffer{}, err
	}
	defer session.Close()

	var buffer bytes.Buffer
	session.Stdout = &buffer
	session.Stdin = os.Stdin
	session.Stderr = os.Stderr

	err = session.Run(fmt.Sprintf("cat %s; rm -f %s", filePath, filePath))
	return buffer, err
}
