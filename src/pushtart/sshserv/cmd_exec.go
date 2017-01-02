package sshserv

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"path"
	"pushtart/config"
	"pushtart/logging"
	"pushtart/tartmanager"
	"pushtart/user"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

func extractPushURL(cmdStr string) string {
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
		err := tartmanager.PreGitRecieve(extractPushURL(cmdStr), conn.User())
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

		sendExitStatus(channel, 0)
	} else if strings.HasPrefix(cmdStr, "git-upload-pack") {
		cmd := exec.Command("git-upload-pack", getPath(cmdStr))
		err := runCommandAcrossSSHChannel(cmd, channel)
		if err != nil {
			logging.Error("sshserv-exec", "runCommandAcrossSSHChannel() returned error: "+err.Error())
			sendExitStatus(channel, 1)
			return
		}
		sendExitStatus(channel, 0)
	} else if strings.HasPrefix(cmdStr, "import-ssh-key ") {
		runImportSSHKey(channel, conn, cmdStr)
	} else if cmdStr == "logs" {
		runLog(channel, conn, cmdStr)
	} else {
		logging.Warning("sshserv-exec", "Exec request disallowed: "+cmdStr)
		channel.Write([]byte("Invalid command - are you using git push?\r\n"))
	}
}

func runLog(channel ssh.Channel, conn *ssh.ServerConn, cmdStr string) {
	bklog := logging.GetBacklog()
	for _, msg := range bklog {
		_, err := fmt.Fprintln(channel, time.Unix(msg.Created, 0).Format(time.ANSIC), "["+msg.Component+"]", msg.Message)
		if err != nil {
			return
		}
	}

	in := make(chan logging.LogMessage, 2)
	logging.Subscribe(in)
	defer logging.Unsubscribe(in)

	for msg := range in {
		_, err := fmt.Fprintln(channel, time.Unix(msg.Created, 0).Format(time.ANSIC), "["+msg.Component+"]", msg.Message)
		if err != nil {
			return
		}
	}
}

func runImportSSHKey(channel ssh.Channel, conn *ssh.ServerConn, cmdStr string) {
	spl := strings.Split(cmdStr, " ")
	if len(spl) < 3 || spl[1] != "--username" {
		channel.Write([]byte("USAGE: import-ssh-key --username <username>\r\n"))
		return
	}

	d, err := ioutil.ReadAll(channel)
	if err != nil {
		logging.Error("sshserv-exec-importssh", "Read err: "+err.Error())
		channel.Write([]byte("ERR: read error." + err.Error() + ".\r\nAbort.\r\n"))
		return
	}

	channel.Write([]byte("Read " + strconv.Itoa(len(d)) + " bytes from input.\r\n"))

	if !user.Exists(spl[2]) {
		logging.Error("sshserv-exec-importssh", "User not found: "+spl[2])
		channel.Write([]byte("ERR: user (" + spl[2] + ") not found.\r\nAbort.\r\n"))
		return
	}

	usr := user.Get(spl[2])
	usr.SSHPubKey = string(d)
	user.Save(spl[2], usr)
	channel.Write([]byte("SSH for " + spl[2] + " saved successfully.\r\n"))
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
