package config

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"pushtart/logging"
	"pushtart/util"
	"strconv"
	"sync"
)

// ErrLockfileExists is returned if a Load() is called for a path which is locked by another process.
var ErrLockfileExists = errors.New("Lockfile exists")
var lockFilePath = ""

func readConfig(fpath string) (*Config, error) {
	if err := getConfigLockStatus(fpath); err != nil {
		return nil, err
	}

	var m = &Config{}

	confF, err := os.Open(fpath)

	if err != nil {
		return nil, errors.New("Failed to open config: " + err.Error())
	}
	defer confF.Close()

	dec := json.NewDecoder(confF)

	if err := dec.Decode(&m); err == io.EOF {
	} else if err != nil {
		return nil, errors.New("Failed to decode config: " + err.Error())
	}
	m.Path = fpath
	return m, lockConfig(fpath)
}

func lockConfig(fpath string) error {
	d := strconv.Itoa(os.Getpid())
	err := ioutil.WriteFile(fpath+".lock", []byte(d), 0700)
	if err != nil {
		logging.Error("config", "Failed to write lockfile: "+err.Error())
	} else {
		lockFilePath = fpath + ".lock"
	}
	return err
}

// UnlockConfig should be called as the process is shutting down to allow other processes to read the configuration.
func UnlockConfig() {
	if lockFilePath != "" {
		err := os.Remove(lockFilePath)
		if err != nil {
			logging.Error("config", "Failed to delete lockfile: "+err.Error())
		}
	}
}

func getConfigLockStatus(fpath string) error {
	fpath = fpath + ".lock"
	var exists bool
	var err error
	if exists, err = util.FileExists(fpath); !exists {
		return nil
	}
	if err != nil {
		return err
	}

	d, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}

	pid, convErr := strconv.Atoi(string(d))
	if convErr != nil {
		logging.Error("config", "Failed to parse contents of lock file - are they integer? Length = "+strconv.Itoa(len(d)))
		return convErr
	}

	if os.Getpid() == pid {
		return nil
	}
	return ErrLockfileExists
}

var writeLock sync.Mutex

func writeConfig() (err error) {
	writeLock.Lock()
	defer writeLock.Unlock()
	//logging.Info("config-write", "Now writing to: "+gConfig.Path)

	data, err := json.MarshalIndent(gConfig, "", "  ")
	if err != nil {
		logging.Info("config-write", "Serialization error: "+err.Error())
		return err
	}

	err = ioutil.WriteFile(gConfig.Path, data, 0755)
	if err != nil {
		logging.Error("config-write", "Error saving configuration: "+err.Error())
		return err
	}
	return nil
}

func loadTLS(keyPath, certPath string) (*tls.Config, error) {
	tlsConfig := new(tls.Config)

	tlsConfig.PreferServerCipherSuites = true
	tlsConfig.CipherSuites = []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA}

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		logging.Error("config", "Error loading tls certificate and key files.")
		logging.Error("config", err.Error())
		return nil, err
	}

	tlsConfig.Certificates = []tls.Certificate{cert}
	tlsConfig.BuildNameToCertificate()

	return tlsConfig, nil
}
