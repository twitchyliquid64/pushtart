package tartmanager

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"pushtart/config"
	"pushtart/logging"
	"pushtart/sshserv/cmd_registry"
	"pushtart/util"
	"strings"
)

// ErrTartOperationNotAuthorized is returned if the pushing user is not allowed to perform that action.
var ErrTartOperationNotAuthorized = errors.New("You are not authorized to perform the requested action.")

func getRepoPath(pushURL string) string {
	return path.Join(config.All().DataPath, pushURL)
}
func getDeploymentPath(pushURL string) string {
	return path.Join(config.All().DeploymentPath, pushURL)
}

func checkCreateRepo(pushURL, owner string) error {
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

		if !checkUserHasTartOwnership(owner, tart.Owners) {
			logging.Warning("tartmanager-git-hooks", "Aborting git-recieve for tart '"+pushURL+"'. Pushing user is not the owner of the tart.")
			return ErrTartOperationNotAuthorized
		}

		if tart.IsRunning {
			err := Stop(pushURL)
			if err != nil {
				logging.Info("tartmanager-git-hooks", "Failed to stop tart: "+err.Error())
			}
		}
		return nil
	}

	return nil
}

// PreGitRecieve is called by the sshserv package when a git push is recieved. It initializes a new repository if one does not already
// exist.
func PreGitRecieve(pushURL, owner string) error {
	return checkCreateRepo(pushURL, owner)
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

	//Check if there is a tartconfig file
	if exists, _ := util.FileExists(path.Join(getDeploymentPath(pushURL), "tartconfig")); exists {
		err = executeCommandFile(path.Join(getDeploymentPath(pushURL), "tartconfig"), pushURL)
		if err != nil {
			logging.Error("tartmanager-git-hooks", "Failed to execute tartconfig: "+err.Error())
			return err
		}
	}

	err = Start(pushURL)
	if err != nil {
		logging.Error("tartmanager-git-hooks", "Failed to start tart: "+err.Error())
	}
	return err
}

func executeCommandFile(fPath, pushURL string) error {
	b, err := ioutil.ReadFile(fPath)
	if err != nil {
		return err
	}

	for _, line := range strings.Split(string(b), "\n") {
		if line == "" {
			continue
		}

		spl := strings.Split(line, " ")
		logging.Info("tartconfig-exec", "["+pushURL+"] "+line)
		if ok, runFunc := cmd_registry.Command(spl[0]); ok {
			cmd := util.ParseCommands(util.TokeniseCommandString(line[len(spl[0]):]))
			if _, ok := cmd["tart"]; !ok {
				cmd["tart"] = pushURL
			}
			runFunc(cmd, &commandOutputRewriter{PushURL: pushURL})
		}
	}
	return nil
}

type commandOutputRewriter struct {
	PushURL string
}

func (c *commandOutputRewriter) Write(p []byte) (n int, err error) {
	logging.Info("tartconfig-exec", "["+c.PushURL+"] "+strings.Replace(string(p), "\n", "", -1))
	return len(p), nil
}

func checkUserHasTartOwnership(user string, owners []string) bool {
	for _, owner := range owners {
		if user == owner {
			return true
		}
	}
	return false
}
