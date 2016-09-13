package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"pushtart/constants"
	"pushtart/logging"
	"strconv"
	"path"
	"pushtart/util"

	"golang.org/x/crypto/ssh"
)

func Generate(fpath string) (err error) {
	logging.Info("config", "Now generating default config to: "+fpath)

	if gConfig == nil {
		gConfig = &Config{
			Name: "pushtart",
			Path: fpath,
		}
	}

	if gConfig.SSH.PrivPEM == "" {
		gConfig.SSH.PubPEM, gConfig.SSH.PrivPEM, err = MakeSSHKeyPair()
	}

	if gConfig.SSH.Listener == "" {
		gConfig.SSH.Listener = "0.0.0.0:2022"
	}

	if gConfig.DataPath == "" {
		pwd, _ := os.Getwd()
		gConfig.DataPath = path.Join(pwd, "gitdata")
		if exists, _ := util.DirExists(gConfig.DataPath); !exists {
			os.Mkdir(gConfig.DataPath, 0777)
		}
	}

	if gConfig.DeploymentPath == "" {
		pwd, _ := os.Getwd()
		gConfig.DeploymentPath = path.Join(pwd, "deploymentdata")
		if exists, _ := util.DirExists(gConfig.DeploymentPath); !exists {
			os.Mkdir(gConfig.DeploymentPath, 0777)
		}
	}


	return writeConfig()
}

func Flush() {
	writeConfig()
}

// MakeSSHKeyPair make a pair of public and private keys for SSH access.
// Public key is encoded in the format for inclusion in an OpenSSH authorized_keys file.
// Private Key generated is PEM encoded
// Source: http://stackoverflow.com/questions/21151714/go-generate-an-ssh-public-key
func MakeSSHKeyPair() (pubKey, privKey string, err error) {
	logging.Info("config-generate", "Now generating SSH private key.")
	logging.Info("config-generate", "Key scheme: RSA. Key size: "+strconv.Itoa(constants.RsaKeySize))
	privateKey, err := rsa.GenerateKey(rand.Reader, constants.RsaKeySize)
	if err != nil {
		return "", "", err
	}

	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	privKey = string(pem.EncodeToMemory(privateKeyPEM))

	// generate and write public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	pubKey = string(ssh.MarshalAuthorizedKey(pub))
	return
}
