package sshserv

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

func passwordCheck(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
	// Should use constant-time compare (or better, salt+hash) in
	// a production setting.
	if c.User() == "testuser" && string(pass) == "tiger" {
		return nil, nil
	}
	return nil, fmt.Errorf("password rejected for %q", c.User())
}
