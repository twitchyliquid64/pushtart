package sshserv

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"
	"os/exec"
	"path"
	"pushtart/config"
	"pushtart/logging"
	"pushtart/tartmanager"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

func extractPushURL(cmdStr string)string {
	return strings.Replace(cmdStr[17:], "'", "", -1)
}

func getPath(cmdStr string) string {
	return path.Join(config.All().DataPath, extractPushURL(cmdStr))
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

	if strings.HasPrefix(cmdStr, "git-receive-pack") {
		err := tartmanager.PreGitRecieve(extractPushURL(cmdStr))
		if err != nil { //err is already logged.
			sendExitStatus(channel, 1)
			return
		}

		cmd := exec.Command("git-receive-pack", getPath(cmdStr))
		err = runCommandAcrossSSHChannel(cmd, channel)
		if err != nil {
			logging.Error("sshserv-exec", "runCommandAcrossSSHChannel() returned error: "+err.Error())
			sendExitStatus(channel, 1)
			return
		}

		err = tartmanager.PostGitRecieve(extractPushURL(cmdStr), conn.User())
		if err != nil { //err is already logged
			sendExitStatus(channel, 1)
			return
		}

		//channel.Write(makeGitMsg("Hello there", false))
		sendExitStatus(channel, 0)
	} else {
		logging.Warning("sshserv-exec", "Exec request disallowed: "+cmdStr)
		channel.Write([]byte("Invalid command - are you using git push?"))
	}
}

func sendExitStatus(channel ssh.Channel, code int) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(code))
	channel.SendRequest("exit-status", false, buf.Bytes())
}

func pipeChannelCopyRoutine(name string, dst io.Writer, src io.Reader, wg *sync.WaitGroup) {
	_, err := io.Copy(dst, src)
	if err != nil {
		logging.Error("sshserv-exec", name+" error: "+err.Error())
	}
	wg.Done()
}

func runCommandAcrossSSHChannel(cmd *exec.Cmd, channel ssh.Channel) error {
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

	wg.Add(1) //copy remote --> command
	go pipeChannelCopyRoutine("stdin", stdinP, channel, &wg)
	wg.Add(1) //copy command --> remote
	go pipeChannelCopyRoutine("stdout", channel, stdout, &wg)
	wg.Add(1) //copy command --> remote (stderr)
	go pipeChannelCopyRoutine("stderr", channel.Stderr(), stderr, &wg)

	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	wg.Wait()
	if err != nil {
		return err
	}
	return nil
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
