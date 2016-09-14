package config

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"pushtart/logging"
	"sync"
)

func readConfig(fpath string) (*Config, error) {
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
	return m, nil
}

var writeLock sync.Mutex

func writeConfig() (err error) {
	writeLock.Lock()
	writeLock.Unlock()

	data, err := json.MarshalIndent(gConfig, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(gConfig.Path, data, 0755)
	if err != nil {
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
