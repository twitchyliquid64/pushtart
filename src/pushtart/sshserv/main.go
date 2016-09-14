package sshserv

import "golang.org/x/crypto/ssh"

var gConfig *ssh.ServerConfig

//Init is called to start the ssh server listener, and have it service requests.
func Init() error {
	return initServConfig()
}
