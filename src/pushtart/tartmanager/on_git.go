package tartmanager

import (
	"os"
	"os/exec"
	"path"
	"pushtart/config"
	"pushtart/logging"
	"pushtart/util"
)

func getRepoPath(pushURL string) string {
	return path.Join(config.All().DataPath, pushURL)
}
func getDeploymentPath(pushURL string) string {
	return path.Join(config.All().DeploymentPath, pushURL)
}

func checkCreateRepo(pushURL string) error {
	repoPath := getRepoPath(pushURL)

	if tartExists := Exists(pushURL); !tartExists {
		logging.Info("tartmanager-git-hooks", "Recieving git push for previously-unknown tart ("+pushURL+").")
		exist, _ := util.DirExists(repoPath)
		if !exist {
			err := os.Mkdir(repoPath, 0777)
			if err != nil {
				logging.Error("tartmanager-git-hooks", "Error creating repository directory: "+err.Error())
				return err
			}
			logging.Info("tartmanager-git-hooks", "Repository directory created.")
		}

		cmd := exec.Command("git", "init", "--bare")
		cmd.Dir = repoPath
		err := cmd.Run()
		if err != nil {
			logging.Error("tartmanager-git-hooks", "Error running git init on repository: "+err.Error())
			return err
		}
		logging.Info("tartmanager-git-hooks", "Repository directory initialized (--bare is SET).")
	} else {
		logging.Info("tartmanager-git-hooks", "Receiving git push for existing tart: "+pushURL)
		tart := Get(pushURL)
		if tart.IsRunning {
			logging.Info("tartmanager-git-hooks", "Tart is currently running. Killing PID ", tart.PID)

			proc, err := os.FindProcess(tart.PID)
			if err != nil {
				return nil
			}
			proc.Kill()
		}
		return nil
	}

	return nil
}

// PreGitRecieve is called by the sshserv package when a git push is recieved. It initializes a new repository if one does not already
// exist.
func PreGitRecieve(pushURL string) error {
	return checkCreateRepo(pushURL)
}

// PostGitRecieve is called after a successful git push. It erases the old deployment if one exists, deploys the new files,
// updates (or creates) the tart object, and finally launches the tart.
func PostGitRecieve(pushURL, owner string) error {
	if !Exists(pushURL) {
		logging.Info("tartmanager-git-hooks", "Registering new tart.")
		New(pushURL, owner)
	} else {
		logging.Info("tartmanager-git-hooks", "Deleting old deployment directory.")
		cmd := exec.Command("rm", "-rf", getDeploymentPath(pushURL))
		cmd.Output()
	}

	cmd := exec.Command("mkdir", getDeploymentPath(pushURL))
	_, err := cmd.Output()
	if err != nil {
		logging.Error("tartmanager-git-hooks", "Failed to create deployment directory: "+err.Error())
		return err
	}

	cmd = exec.Command("git", "clone", getRepoPath(pushURL), "./")
	cmd.Dir = getDeploymentPath(pushURL)
	_, err = cmd.Output()
	if err != nil {
		logging.Error("tartmanager-git-hooks", "Failed to clone repository to deployment directory: "+err.Error())
		return err
	}
	return nil
}
