package config

import (
	"crypto/tls"
	"pushtart/logging"
)

var gConfig *Config = nil
var gTls *tls.Config = nil

func Load(fpath string) error {
	conf, err := readConfig(fpath)
	if err == nil {
		gConfig = conf
	} else {
		logging.Error("config", "config.Load() error:", err)
		return err
	}

	if gConfig.TLS.PrivateKey == "" {
		logging.Warning("config", "TLS keyfile paths omitted, skipping TLS setup")
	} else {
		tls, err := loadTLS(gConfig.TLS.PrivateKey, gConfig.TLS.Cert)
		if err == nil {
			gTls = tls
		} else {
			logging.Error("config", "config.Load() tls error:", err)
			return err
		}
	}

	return nil
}

func GetServerName() string {
	checkInitialisedOrPanic()
	return gConfig.Name
}

func TLS() *tls.Config {
	checkInitialisedOrPanic()
	return gTls
}

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
