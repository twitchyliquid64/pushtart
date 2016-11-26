package config

import (
	"crypto/tls"
	"errors"
	"pushtart/logging"
)

var gConfig *Config
var gTLS *tls.Config

//Load loads the configuration JSON file at fpath as global configuration.
func Load(fpath string) error {
	conf, err := readConfig(fpath)
	if err == nil {
		gConfig = conf
	} else {
		logging.Error("config", "config.Load() error: ", err)
		return err
	}

	if gConfig.TLS.PrivateKey == "" {
		//logging.Warning("config", "TLS keyfile paths omitted, skipping TLS setup")
	} else {
		tls, err := loadTLS(gConfig.TLS.PrivateKey, gConfig.TLS.Cert)
		if err == nil {
			gTLS = tls
		} else {
			logging.Error("config", "config.Load() tls error:", err)
			return err
		}
	}

	if gConfig.RunSentryInterval == 0 {
		logging.Error("config", "RunSentryInterval cannot be 0")
		return errors.New("RunSentryInterval cannot be 0")
	}

	return nil
}

//GetServerName returns the field 'Name' specified in the configuration.
func GetServerName() string {
	checkInitialisedOrPanic()
	return gConfig.Name
}

//TLS returns a TLS configuration object used to setup https servers.
func TLS() *tls.Config {
	checkInitialisedOrPanic()
	return gTLS
}

//All returns the full configuration object - use this to access arbitrary fields.
func All() *Config {
	checkInitialisedOrPanic()
	return gConfig
}

func checkInitialisedOrPanic() {
	if gConfig == nil {
		panic("Config not initialised")
	}
	//if gTls == nil{
	//	panic("TLS not initialised")
	//}
}
