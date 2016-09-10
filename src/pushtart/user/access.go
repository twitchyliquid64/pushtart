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

func New(username string){
  if config.All().Users == nil{
    config.All().Users = map[string]config.User{}
  }

  config.All().Users[username] = config.User{
  }
  config.Flush()
}
