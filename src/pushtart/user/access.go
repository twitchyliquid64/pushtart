package user

import (
  "pushtart/config"
)

func Get(username string)config.User{
  return config.All().Users[username]
}


func Exists(username string)bool{
  _, ok := config.All().Users[username]
  return ok
}

func Save(username string, usr config.User){
  if config.All().Users == nil{
    config.All().Users = map[string]config.User{}
  }
  config.All().Users[username] = usr
  config.Flush()
}

func New(username string){
  Save(username,config.User{})
}
