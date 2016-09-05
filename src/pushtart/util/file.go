package util

import (
  "pushtart/logging"
  "strings"
  "io/ioutil"
  "path"
  "os"
)


func GetFilenameListInFolder(folder, suffix string)([]string,error){
  output := []string{}

  pwd, err := os.Getwd()
    if err != nil {
        logging.Error("file-util", err)
        return nil, err
    }

    files, err := ioutil.ReadDir(path.Join(pwd, folder))
  	if err != nil {
      logging.Error("file-util", err)
      return nil, err
  	}

    for _, file := range files {
      if (!file.IsDir()) && strings.HasSuffix(file.Name(), suffix){

        p := file.Name()
        if !path.IsAbs(p){
          p = path.Join(path.Join(pwd, folder), file.Name())
        }
        output = append(output, p)
      }
    }
    return output, nil
}
