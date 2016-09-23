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

	DNS struct {
		Enabled         bool
		Listener        string
		AllowForwarding bool
		LookupCacheSize int
		ARecord         map[string]ARecord
		AAAARecord      map[string]ARecord
	}

	Users map[string]User
	Tarts map[string]Tart
}

//ARecord represents a response that could be served to a DNS query of type A.
type ARecord struct {
	Address string
	TTL     uint32
}

//User represents an account which has access to the system.
type User struct {
	Name             string
	Password         string
	AllowSSHPassword bool
	SSHPubKey        string
}

//Tart stores information for tarts which are stored in the system.
type Tart struct {
	PushURL          string
	Name             string
	Owner            string
	IsRunning        bool
	LogStdout        bool
	PID              int
	Env              []string
	RestartOnStop    bool
	RestartDelaySecs int
}
