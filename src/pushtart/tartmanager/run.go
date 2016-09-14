package tartmanager

import (
	"errors"
	"os"
	"os/exec"
	"path"
	"pushtart/logging"
	"pushtart/util"
	"strings"
)

//ErrTartNotFound is returned if a stop/start is requested but the tart specified by pushURL could not be found.
var ErrTartNotFound = errors.New("Tart not found")

//ErrTartWrongState is returned if a stop is requested on a stopped tart, or a start is requested on a running tart.
var ErrTartWrongState = errors.New("Tart is in the wrong state to execute that command.")

const runScriptSh = "startup.sh"
const runScriptPy = "startup.py"

//Start commences execution of the given tart.
func Start(pushURL string) error {
	if !Exists(pushURL) {
		return ErrTartNotFound
	}
	tart := Get(pushURL)
	if tart.IsRunning {
		return ErrTartWrongState
	}
	deploymentFolder := getDeploymentPath(pushURL)

	var cmd *exec.Cmd
	if shExists, _ := util.FileExists(path.Join(deploymentFolder, runScriptSh)); shExists {
		cmd = exec.Command("bash", runScriptSh)
	} else if pyExists, _ := util.FileExists(path.Join(deploymentFolder, runScriptPy)); pyExists {
		cmd = exec.Command("python", runScriptPy)
	} else {
		return errors.New("No startup script")
	}

	cmd.Dir = deploymentFolder
	err := cmd.Start()
	if err != nil {
		return err
	}

	tart.PID = cmd.Process.Pid
	tart.IsRunning = true
	Save(pushURL, tart)
	logging.Info("tartmanager-run", "Started "+pushURL)
	return nil
}

//Stop halts execution of the given tart.
func Stop(pushURL string) error {
	if !Exists(pushURL) {
		return ErrTartNotFound
	}
	tart := Get(pushURL)
	if !tart.IsRunning {
		return ErrTartWrongState
	}

	logging.Info("tartmanager-run", "Killing running tart with PID ", tart.PID)

	proc, err := os.FindProcess(tart.PID)
	if err != nil {
		if strings.Contains(err.Error(), "process already finished") {
			logging.Warning("tartmanager-run", "Aborting stop operation on "+pushURL+", process already terminated.")
		} else {
			return err
		}
	}
	tart.PID = -1
	tart.IsRunning = false
	Save(pushURL, tart)
	err = proc.Kill()

	if err != nil && strings.Contains(err.Error(), "process already finished") {
		logging.Warning("tartmanager-run", "Aborting stop operation on "+pushURL+", process already terminated.")
		return nil
	}
	return err
}
