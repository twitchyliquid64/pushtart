package sshserv

import (
	"pushtart/config"

	"golang.org/x/crypto/ssh"
)

func initServConfig() (err error) {
	gConfig = &ssh.ServerConfig{
		PasswordCallback:  passwordCheck,
		PublicKeyCallback: publicKeyCheck,
	}

	private, err := ssh.ParsePrivateKey([]byte(config.All().SSH.PrivPEM))
	if err != nil {
		return err
	}
	gConfig.AddHostKey(private)

	return startListener()
}
