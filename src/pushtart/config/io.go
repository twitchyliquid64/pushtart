package config

import (
	"pushtart/logging"
	"encoding/json"
	"crypto/tls"
	"errors"
	"io"
	"os"
)

func readConfig(fpath string)(*Config,error){
	var m = &Config{}

	confF, err := os.Open(fpath)

	if err != nil {
			return nil, errors.New("Failed to open config: "+err.Error())
	}
	defer confF.Close()

	dec := json.NewDecoder(confF)

	if err := dec.Decode(&m); err == io.EOF {
	} else if err != nil {
		return nil, errors.New("Failed to decode config: "+err.Error())
	}
	return m, nil
}

func loadTLS(keyPath, certPath string)(*tls.Config, error){
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
