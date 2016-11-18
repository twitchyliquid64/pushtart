package user

import (
	"pushtart/config"
	"pushtart/logging"
	"pushtart/util"

	"golang.org/x/crypto/bcrypt"
)

//GetUserPubkey returns the openssh-friendly user public key.
func GetUserPubkey(username string) string {
	usrStruct, ok := config.All().Users[username]
	if !ok {
		return ""
	}

	return usrStruct.SSHPubKey
}

// CheckUserPasswordSSH checks the given username to see if the stored password matches the one provided.
// If the user does not exist, or the passwords do not match, false is returned. If the password matches
// and AllowSSHPassword is true, the function returns true.
func CheckUserPasswordSSH(username, password string) bool {
	usrStruct, ok := config.All().Users[username]
	if !ok {
		return false
	}

	if usrStruct.AllowSSHPassword {
		err := util.ComparePassHash(usrStruct.Password, username, password)
		if err == nil { //if err == nil the passwords match
			return true
		}

		if err != bcrypt.ErrMismatchedHashAndPassword {
			logging.Error("sshpwd-auth", "Hash compare error: "+err.Error())
		}
	} else {
		logging.Error("sshpwd-auth", "Password authentication attempted for user with AllowSSHPassword = false. (Have you run edit-user with '--allow-ssh-password yes'?)")
	}
	return false
}

// CheckUserPasswordWeb checks the given username to see if the stored password matches the one provided.
// If the user does not exist, or the passwords do not match, false is returned.
func CheckUserPasswordWeb(username, password string) bool {
	usrStruct, ok := config.All().Users[username]
	if !ok {
		return false
	}

	// Speedup for web clients making heaps of requests - bcrypt is expensive
	if isValidFromCacheHit(username, password) {
		return true
	}

	err := util.ComparePassHash(usrStruct.Password, username, password)
	if err == nil { //if err == nil the passwords match
		cacheCorrectAuthEntry(username, password)
		return true
	}

	if err != bcrypt.ErrMismatchedHashAndPassword {
		logging.Error("httpproxy-auth", "Hash compare error: "+err.Error())
	}

	return false
}
