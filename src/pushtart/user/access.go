package user

import (
	"errors"
	"pushtart/config"
)

//Get returns the user structure for the given useername.
func Get(username string) config.User {
	return config.All().Users[username]
}

//Exists returns true if the given username exists.
func Exists(username string) bool {
	_, ok := config.All().Users[username]
	return ok
}

//Save writes the given user structure to global configuration under the given username, before flushing the global configuration to disk.
func Save(username string, usr config.User) {
	if config.All().Users == nil {
		config.All().Users = map[string]config.User{}
	}
	config.All().Users[username] = usr
	config.Flush()
}

//Delete removes the specified user from the system.
func Delete(username string) error {
	if config.All().Users == nil {
		config.All().Users = map[string]config.User{}
	}

	for _, tart := range config.All().Tarts {
		for _, owner := range tart.Owners {
			if owner == username {
				return errors.New("Cannot delete a user who still owns a tart")
			}
		}
	}

	delete(config.All().Users, username)
	config.Flush()
	return nil
}

//New creates a new user wit the given username, writing it to global configuration before flushing to disk.
func New(username string) {
	Save(username, config.User{})
}

//List returns a []string of all the users in the system.
func List() []string {
	var output []string
	for username := range config.All().Users {
		output = append(output, username)
	}
	return output
}
