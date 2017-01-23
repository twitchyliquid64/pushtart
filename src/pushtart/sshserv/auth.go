package sshserv

import (
	"bytes"
	"errors"
	"fmt"
	"pushtart/logging"
	"pushtart/user"

	"golang.org/x/crypto/ssh"
)

var errAuthDenied = errors.New("Authentication denied")

func passwordCheck(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
	logging.Info("sshserv-auth", "Received authentication request for "+c.User()+" (Password)")

	// Should use constant-time compare (or better, salt+hash) in
	// a production setting.
	if user.CheckUserPasswordSSH(c.User(), string(pass)) {
		return nil, nil
	}
	return nil, fmt.Errorf("password rejected for %q", c.User())
}

func publicKeyCheck(conn ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
	logging.Info("sshserv-auth", "Received authentication request for "+conn.User()+" (PublicKey)")

	pubKeyRaw := user.GetUserPubkey(conn.User())
	if pubKeyRaw == "" {
		logging.Warning("sshserv-auth", "No PublicKey known for "+conn.User())
		return nil, errAuthDenied
	}

	trustedPk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(pubKeyRaw))
	if err != nil {
		logging.Error("sshserv-auth", "Could not parse user's public key: "+err.Error())
		return nil, errAuthDenied
	}

	if bytes.Equal(trustedPk.Marshal(), pubKey.Marshal()) { //authentication successful
		return nil, nil
	}
	return nil, errAuthDenied
}
