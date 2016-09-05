package sshserv

import "golang.org/x/crypto/ssh"

var gConfig *ssh.ServerConfig

func Init() {
	initServConfig()
}
