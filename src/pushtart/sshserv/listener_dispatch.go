package sshserv

import (
	"errors"
	"net"
	"pushtart/config"
	"pushtart/logging"

	"golang.org/x/crypto/ssh"
)

func startListener() error {
	if config.All().Ssh.Listener == "" {
		return errors.New("No listener specified in config")
	}

	listener, err := net.Listen("tcp", config.All().Ssh.Listener)
	if err != nil {
		return err
	}

	for { // accept + handshake routine
		newConn, err := listener.Accept()
		if err != nil {
			logging.Error("sshserv-accept", err.Error())
			continue
		}
		go handshakeSocket(newConn)
	}
	return nil
}

func handshakeSocket(newConn net.Conn) {
	sshConn, chans, reqs, err := ssh.NewServerConn(newConn, gConfig)
	if err != nil {
		logging.Error("sshserv-handshake", err.Error()+": "+newConn.RemoteAddr().String())
		return
	}
	go ssh.DiscardRequests(reqs) // The incoming Request channel must be serviced.
	serviceSshConnection(sshConn, chans)
}

func serviceSshConnection(conn *ssh.ServerConn, newSshChannelReq <-chan ssh.NewChannel) {
	for newChannel := range newSshChannelReq {
		// Channels have a type, depending on the application level protocol intended. In the case of a shell, the type is
		// "session" and ServerShell may be used to present a simple terminal interface.
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		logging.Info("sshserv-service", "Got channel request: "+newChannel.ChannelType()+" ("+conn.User()+")")

		_, _, err := newChannel.Accept()
		if err != nil {
			logging.Error("sshserv-service", err.Error()+": "+conn.RemoteAddr().String())
		}
	}
}
