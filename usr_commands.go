package main

import (
  "strings"
  "pushtart/user"
  "pushtart/util"
  "pushtart/config"
  "pushtart/logging"
  "io/ioutil"
  "os"
  "io"
  "bytes"
  "fmt"
)

func saveUser(username string, usr config.User, params map[string]string){
	var exists bool

	if _, exists = params["password"]; exists {
		pw, err := util.HashPassword(username, params["password"])
		if err != nil {
			logging.Error("make-user", "Error hashing password: " + err.Error())
			return
		}
		usr.Password = pw
	}

	if _, exists = params["name"]; exists {
		usr.Name = params["name"]
	}

	if _, exists = params["allow-ssh-password"]; exists {
		usr.AllowSSHPassword = false
		if strings.ToUpper(params["allow-ssh-password"]) == "YES" {
			usr.AllowSSHPassword = true
		}
	}

	user.Save(username, usr)
}


func editUser(params map[string]string){
	if missingFields := checkHasFields([]string{"username"}, params); len(missingFields) > 0 {
		fmt.Println("USAGE: pushtart edit-user --username <username> [--config <config file>] [--password <password] [--name <name] [--allow-ssh-password yes/no]")
		printMissingFields(missingFields)
		return
	}

  if !user.Exists(params["username"]){
    fmt.Println("Err: user does not exist")
    return
  }

	usr := user.Get(params["username"])
	saveUser(params["username"], usr, params)
}

func makeUser(params map[string]string){
	if missingFields := checkHasFields([]string{"username"}, params); len(missingFields) > 0 {
		fmt.Println("USAGE: pushtart make-user --username <username> [--config <config file>] [--password <password] [--name <name] [--allow-ssh-password yes/no]")
		printMissingFields(missingFields)
		return
	}

	if user.Exists(params["username"]){
		fmt.Println("Err: user already exists")
		return
	}

	user.New(params["username"])
	usr := user.Get(params["username"])
	saveUser(params["username"], usr, params)
}


func importSshKey(params map[string]string){
  if missingFields := checkHasFields([]string{"username"}, params); len(missingFields) > 0 {
		fmt.Println("USAGE: pushtart import-ssh-key --username <username> [--pub-key-file <path-to-.pub-file>]")
		printMissingFields(missingFields)
		return
	}

  if !user.Exists(params["username"]){
    fmt.Println("Err: user does not exist")
    return
  }

  var err error
  var b []byte
  if _, pathExists := params["pub-key-file"]; pathExists {
    b, err = ioutil.ReadFile(params["pub-key-file"])
  } else {
    buf := bytes.NewBuffer(nil)
    _, err = io.Copy(buf, os.Stdin)
    b = buf.Bytes()
  }

  if err != nil{
    logging.Error("import-ssh-key", "Read error: " + err.Error())
    return
  }

  usr := user.Get(params["username"])
  usr.SSHPubKey = string(b)
  user.Save(params["username"], usr)
}
