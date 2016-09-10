package sshserv

import (
	"errors"
	"net"
	"pushtart/config"
	"pushtart/logging"

	"golang.org/x/crypto/ssh"
)

func startListener() error {
	if config.All().SSH.Listener == "" {
		return errors.New("No listener specified in config")
	}

	listener, err := net.Listen("tcp", config.All().SSH.Listener)
	if err != nil {
		return err
	}

	go func() {
		for { // accept routine
			newConn, err := listener.Accept()
			if err != nil {
				logging.Error("sshserv-accept", err.Error())
				continue
			}
			go handshakeSocket(newConn)
		}
	}()
	return nil
}

func handshakeSocket(newConn net.Conn) {
	sshConn, chans, reqs, err := ssh.NewServerConn(newConn, gConfig)
	if err != nil {
		logging.Error("sshserv-handshake", err.Error()+": "+newConn.RemoteAddr().String())
		return
	}
	go ssh.DiscardRequests(reqs) // The incoming Request channel must be serviced.
	serviceSSHConnection(sshConn, chans)
}

func serviceSSHConnection(conn *ssh.ServerConn, newSSHChannelReq <-chan ssh.NewChannel) {
	for newChannel := range newSSHChannelReq {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		//logging.Info("sshserv-service", "Got channel request: "+newChannel.ChannelType()+" ("+conn.User()+")")

		channel, requests, err := newChannel.Accept()
		if err != nil {
			logging.Error("sshserv-service", err.Error()+": "+conn.RemoteAddr().String())
			continue
		}

		go serviceSSHChannel(conn, channel, requests)
	}
}

func serviceSSHChannel(conn *ssh.ServerConn, channel ssh.Channel, requests <-chan *ssh.Request) {
	for req := range requests {
		//logging.Info("sshserv-service", "Got OOB request ("+string(req.Type)+")")
		switch req.Type {

		case "shell":
			req.Reply(true, nil)
			go shell(conn, channel)

		case "pty-req": //need this to get shell to work for some reason
			req.Reply(true, nil)

		case "exec":
			req.Reply(true, nil)
			go execCmd(conn, channel, req.Payload)

		default:
			req.Reply(false, nil) //err
		}
	}
}
