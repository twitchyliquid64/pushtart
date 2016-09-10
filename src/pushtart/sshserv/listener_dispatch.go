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

	go func(){
		for { // accept + handshake routine
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

		channel, requests, err := newChannel.Accept()
		if err != nil {
			logging.Error("sshserv-service", err.Error()+": "+conn.RemoteAddr().String())
		}

		go serviceSshChannel(conn, channel, requests)
	}
}


func serviceSshChannel(conn *ssh.ServerConn, channel ssh.Channel, requests <-chan *ssh.Request){
	for req := range requests {
			logging.Info("sshserv-service", "Got OOB request (" + string(req.Type) + ")")
			switch req.Type {
			case "shell":
				req.Reply(true, nil)
				go shell(conn, channel)
			default:
				req.Reply(false, nil)//err
			}
	}
	logging.Info("sshserv-service", "Channel closing")
}
