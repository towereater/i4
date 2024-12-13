package utils

import (
	"aggregator/config"
	"aggregator/model"
	"fmt"
	"os"
	"path"

	"golang.org/x/crypto/ssh"
)

func GetFileFromRemote(cfg config.Config, t model.Target) (*os.File, error) {
	clientConfig := &ssh.ClientConfig{
		User: t.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(t.Pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", t.NetIp, clientConfig)
	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	f, err := CreateOrReplaceFile(path.Join(cfg.FileDir, fmt.Sprintf("%s-%s", t.Id, t.Machine)))
	if err != nil {
		return nil, err
	}

	session.Stdout = f
	session.Stdin = os.Stdin
	session.Stderr = os.Stderr

	err = session.Run(fmt.Sprintf("cat %s; rm -f %s", t.File, t.File))
	return f, err
}
