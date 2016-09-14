package config

//Config represents the global configuration of the system, and the JSON structure on disk.
type Config struct {
	Name           string //canonical name to help identify this server
	DataPath       string
	DeploymentPath string
	Path           string   `json:"-"` //path used to represent where the file is currently stored.
	TLS            struct { //Relative file addresses of the .pem files needed for TLS.
		PrivateKey string
		Cert       string
	}

	Web struct { //Details needed to get the website part working.
		Domain   string //Domain should be in the form example.com
		Listener string //Address:port (address can be omitted) where the HTTPS listener will bind.
	}

	SSH struct {
		PubPEM   string
		PrivPEM  string
		Listener string
	}

	Users map[string]User
}

//User represents an account which has access to the system.
type User struct {
	Name             string
	Password         string
	AllowSSHPassword bool
	SSHPubKey        string
}
