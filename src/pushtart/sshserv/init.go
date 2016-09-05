package sshserv

import (
	"pushtart/config"

	"golang.org/x/crypto/ssh"
)

func initServConfig() (err error) {
	gConfig = &ssh.ServerConfig{
		PasswordCallback: passwordCheck,
	}

	private, err := ssh.ParsePrivateKey([]byte(config.All().Ssh.PrivPEM))
	if err != nil {
		return err
	}
	gConfig.AddHostKey(private)

	return nil
}
