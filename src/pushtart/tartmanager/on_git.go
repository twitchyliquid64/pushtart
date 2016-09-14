package tartmanager

import (
  "pushtart/util"
  "pushtart/config"
  "pushtart/logging"
  "os"
  "path"
  "os/exec"
)

func getPath(pushURL string) string {
	return path.Join(config.All().DataPath, pushURL)
}

func checkCreateRepo(pushURL string)error {
  repoPath := getPath(pushURL)

	if tartExists := Exists(pushURL); !tartExists{
		logging.Info("tartmanager-git-hooks", "Recieving git push for previously-unknown tart (" + pushURL + ").")
    exist, _ := util.DirExists(repoPath)
  	if !exist {
  		err := os.Mkdir(repoPath, 0777)
      if err != nil {
        logging.Error("tartmanager-git-hooks", "Error creating repository directory: " + err.Error())
        return err
      } else {
        logging.Info("tartmanager-git-hooks", "Repository directory created.")
      }
  	}

  	cmd := exec.Command("git", "init", "--bare")
  	cmd.Dir = repoPath
  	err := cmd.Run()
  	if err != nil {
  		logging.Error("tartmanager-git-hooks", "Error running git init on repository: "+err.Error())
      return err
  	} else {
      logging.Info("tartmanager-git-hooks", "Repository directory initialized (--bare is SET).")
    }
  } else {
    logging.Info("tartmanager-git-hooks", "Receiving git push for existing tart: " + pushURL)
    tart := Get(pushURL)
    if tart.IsRunning {
      logging.Info("tartmanager-git-hooks", "Tart is currently running. Killing PID ", tart.PID)

      proc, err := os.FindProcess(tart.PID)
      if err != nil{
        return
      }
      proc.Kill()
    }
    return nil
  }

  return nil
}



func PreGitRecieve(pushURL string)error{
  return checkCreateRepo(pushURL)
}

func PostGitRecieve(pushURL, owner string)error {
  if !Exists(pushURL) {
    logging.Error("tartmanager-git-hooks", "Registering new tart.")
    New(pushURL, owner)
  }
  return nil
}
