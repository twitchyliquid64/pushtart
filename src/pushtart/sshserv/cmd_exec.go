package sshserv

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"
	"os"
	"os/exec"
	"path"
	"pushtart/config"
	"pushtart/logging"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

func getPath(cmdStr string) string {
	return path.Join(config.All().DataPath, strings.Replace(cmdStr[17:], "'", "", -1))
}

func checkRepo(cmdStr string) {
	repoPath := getPath(cmdStr)
	exist, _ := dirExists(repoPath)
	if !exist {
		os.Mkdir(repoPath, 0777)
	}
	cmd := exec.Command("git", "init", "--bare")
	cmd.Dir = repoPath
	err := cmd.Run()
	if err != nil {
		logging.Error("sshserv-exec", "pre-init error: "+err.Error())
	}
}

func execCmd(conn *ssh.ServerConn, channel ssh.Channel, payload []byte) {
	cmdStr := string(payload[4:])
	defer func() {
		err := channel.Close()
		if err != nil {
			logging.Error("sshserv-exec", "Close error: "+err.Error())
		}
		logging.Info("sshserv-exec", "Channel closing: "+cmdStr)
	}()

	if !strings.HasPrefix(cmdStr, "git-receive-pack") {
		logging.Warning("sshserv-exec", "Exec request disallowed: "+cmdStr)
		channel.Write([]byte("Invalid command - are you using git push?"))
	} else {
		checkRepo(cmdStr)
		cmd := exec.Command("git-receive-pack", getPath(cmdStr))
		runCommandAcrossSSHChannel(cmd, channel)
		//channel.Write(makeGitMsg("Hello there", false))
	}
}

func pipeChannelCopyRoutine(name string, dst io.Writer, src io.Reader, wg *sync.WaitGroup) {
	_, err := io.Copy(dst, src)
	if err != nil {
		logging.Error("sshserv-exec", name+" error: "+err.Error())
	}
	wg.Done()
}

func runCommandAcrossSSHChannel(cmd *exec.Cmd, channel ssh.Channel) {
	var wg sync.WaitGroup
	stdinP, err := cmd.StdinPipe()
	if err != nil {
		logging.Error("sshserv-exec", "Could not open git-recieve-pack stdin: "+err.Error())
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logging.Error("sshserv-exec", "Could not open git-recieve-pack stdout: "+err.Error())
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		logging.Error("sshserv-exec", "Could not open git-recieve-pack stderr: "+err.Error())
	}

	//copy remote --> command
	wg.Add(1)
	go pipeChannelCopyRoutine("stdin", stdinP, channel, &wg)
	//copy command --> remote
	wg.Add(1)
	go pipeChannelCopyRoutine("stdout", channel, stdout, &wg)
	//copy command --> remote (stderr)
	wg.Add(1)
	go pipeChannelCopyRoutine("stderr", channel.Stderr(), stderr, &wg)

	err = cmd.Start()
	if err != nil {
		logging.Error("sshserv-exec", "cmd.Start() error: "+err.Error())
	}
	err = cmd.Wait()
	if err != nil {
		logging.Error("sshserv-exec", "cmd.Wait() error: "+err.Error())
	}
	wg.Wait()
}

func makeGitMsg(msg string, isError bool) []byte {
	output := ""

	headerBuf := new(bytes.Buffer)
	binary.Write(headerBuf, binary.BigEndian, uint16(len(msg)+4+1))
	output += strings.ToUpper(hex.EncodeToString(headerBuf.Bytes()))

	if isError {
		output += string([]byte{byte(1)})
	} else {
		output += string([]byte{byte(2)})
	}

	output += msg
	return []byte(output)
}

func dirExists(path string) (bool, error) {
	s, err := os.Stat(path)
	if err == nil {
		if s.IsDir() {
			return true, nil
		}
	}
	return false, err
}
