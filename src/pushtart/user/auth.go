package user

import "pushtart/config"

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

	if usrStruct.AllowSSHPassword && usrStruct.Password == password {
		return true
	}
	return false
}
