package sshserv

import (
  "golang.org/x/crypto/ssh"
  "golang.org/x/crypto/ssh/terminal"
  "pushtart/logging"
  "fmt"
)

func shell(conn *ssh.ServerConn, channel ssh.Channel){
  term := terminal.NewTerminal(channel, "> ")

  go func() {
      defer channel.Close()
      for {
          line, err := term.ReadLine()
          if err != nil {
            logging.Error("sshserv-shell", err.Error())
            break
          }
          logging.Info("sshserv-shell", line)
          fmt.Println(line)
      }
  }()
}
