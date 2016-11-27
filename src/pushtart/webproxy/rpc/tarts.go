package rpc

import (
	"errors"
	"pushtart/config"
	"pushtart/logging"
	"pushtart/tartmanager"
	"strconv"
	"strings"

	"github.com/powerman/rpc-codec/jsonrpc2"
)

// Tarts represents the authenticated RPC server for tart-related methods, available via a webproxy URI.
type Tarts int

// GetTartArgument represents the parameters passed to the GetTart RPC.
type GetTartArgument struct {
	APIKey  string
	PushURL string
}

// GetTartResult represents the result of a successful GetTart RPC.
type GetTartResult struct {
	Tart config.Tart
}

// GetTart RPC returns a specific tart
func (t *Tarts) GetTart(arg *GetTartArgument, result *GetTartResult) error {
	var serviceName string
	var ok bool
	if serviceName, ok = checkAuth(arg.APIKey); ok {
		logging.Info("rpc", "["+serviceName+"] GetTart("+arg.PushURL+")")
	} else {
		logging.Warning("rpc", "Invalid auth for GetTart("+arg.PushURL+")")
		return jsonrpc2.NewError(403, "Invalid API key")
	}

	if tartmanager.Exists(arg.PushURL) {
		result.Tart = tartmanager.Get(arg.PushURL)
	} else {
		return errors.New("Could not find tart")
	}

	return nil
}

// GetTartStats RPC returns a running tarts system stats.
func (t *Tarts) GetTartStats(arg *GetTartArgument, result *tartmanager.RunMetrics) error {
	var serviceName string
	var ok bool
	if serviceName, ok = checkAuth(arg.APIKey); ok {
		logging.Info("rpc", "["+serviceName+"] GetTart("+arg.PushURL+")")
	} else {
		logging.Warning("rpc", "Invalid auth for GetTart("+arg.PushURL+")")
		return jsonrpc2.NewError(403, "Invalid API key")
	}

	if tartmanager.Exists(arg.PushURL) {
		res, err := tartmanager.GetStats(arg.PushURL)
		if err != nil {
			return err
		}
		result.Mem = res.Mem
		result.Time = res.Time
		result.State = res.State
		result.Children = res.Children
		result.PID = res.PID
	} else {
		return errors.New("Could not find tart")
	}

	return nil
}

// ArbitrarySuccessResult represents the result of a successful RPC, where only success needs to be indicated.
type ArbitrarySuccessResult struct {
	Success bool
}

// EnableOutputLogging RPC enables Stdout/Error logging to the main logger for a given tart.
func (t *Tarts) EnableOutputLogging(arg map[string]string, result *ArbitrarySuccessResult) error {
	var serviceName string
	var ok bool
	if serviceName, ok = checkAuth(arg["APIKey"]); ok {
		logging.Info("rpc", "["+serviceName+"] EnableOutputLogging("+arg["PushURL"]+")")
	} else {
		logging.Warning("rpc", "Invalid auth for EnableOutputLogging("+arg["PushURL"]+")")
		return jsonrpc2.NewError(403, "Invalid API key")
	}

	if tartmanager.Exists(arg["PushURL"]) {
		t := tartmanager.Get(arg["PushURL"])
		var err error
		t.LogStdout, err = strconv.ParseBool(arg["Enable"])
		if err != nil {
			return err
		}
		tartmanager.Save(arg["PushURL"], t)
		result.Success = true
	} else {
		return errors.New("Could not find tart")
	}

	return nil
}

// SetName RPC sets the human-readable name for a given tart.
func (t *Tarts) SetName(arg map[string]string, result *ArbitrarySuccessResult) error {
	var serviceName string
	var ok bool
	if serviceName, ok = checkAuth(arg["APIKey"]); ok {
		logging.Info("rpc", "["+serviceName+"] SetName("+arg["PushURL"]+", "+arg["Name"]+")")
	} else {
		logging.Warning("rpc", "Invalid auth for SetName("+arg["PushURL"]+")")
		return jsonrpc2.NewError(403, "Invalid API key")
	}

	if tartmanager.Exists(arg["PushURL"]) {
		t := tartmanager.Get(arg["PushURL"])
		t.Name = arg["Name"]
		tartmanager.Save(arg["PushURL"], t)
		result.Success = true
	} else {
		return errors.New("Could not find tart")
	}

	return nil
}

// SetEnv RPC sets a key=value environment variable for the tart.
func (t *Tarts) SetEnv(arg map[string]string, result *ArbitrarySuccessResult) error {
	var serviceName string
	var ok bool
	if serviceName, ok = checkAuth(arg["APIKey"]); ok {
		logging.Info("rpc", "["+serviceName+"] SetEnv("+arg["PushURL"]+", "+arg["Key"]+")")
	} else {
		logging.Warning("rpc", "Invalid auth for SetEnv("+arg["PushURL"]+")")
		return jsonrpc2.NewError(403, "Invalid API key")
	}

	if tartmanager.Exists(arg["PushURL"]) {
		t := tartmanager.Get(arg["PushURL"])
		t.Env = setEnv(t.Env, arg["Key"]+"="+arg["Value"], "")
		tartmanager.Save(arg["PushURL"], t)
		result.Success = true
	} else {
		return errors.New("Could not find tart")
	}

	return nil
}

// DelEnv RPC deletes a key from a tarts environment variable list if it exists.
func (t *Tarts) DelEnv(arg map[string]string, result *ArbitrarySuccessResult) error {
	var serviceName string
	var ok bool
	if serviceName, ok = checkAuth(arg["APIKey"]); ok {
		logging.Info("rpc", "["+serviceName+"] DelEnv("+arg["PushURL"]+", "+arg["Key"]+")")
	} else {
		logging.Warning("rpc", "Invalid auth for DelEnv("+arg["PushURL"]+")")
		return jsonrpc2.NewError(403, "Invalid API key")
	}

	if tartmanager.Exists(arg["PushURL"]) {
		t := tartmanager.Get(arg["PushURL"])
		t.Env = setEnv(t.Env, "", arg["Key"])
		tartmanager.Save(arg["PushURL"], t)
		result.Success = true
	} else {
		return errors.New("Could not find tart")
	}

	return nil
}

// Start RPC starts a tart.
func (t *Tarts) Start(arg map[string]string, result *ArbitrarySuccessResult) error {
	var serviceName string
	var ok bool
	if serviceName, ok = checkAuth(arg["APIKey"]); ok {
		logging.Info("rpc", "["+serviceName+"] Start("+arg["PushURL"]+")")
	} else {
		logging.Warning("rpc", "Invalid auth for Start("+arg["PushURL"]+")")
		return jsonrpc2.NewError(403, "Invalid API key")
	}

	if tartmanager.Exists(arg["PushURL"]) {
		err := tartmanager.Start(arg["PushURL"])
		if err != nil {
			return err
		}
		result.Success = true
		return nil
	}
	return errors.New("Could not find tart")
}

// Stop RPC starts a tart.
func (t *Tarts) Stop(arg map[string]string, result *ArbitrarySuccessResult) error {
	var serviceName string
	var ok bool
	if serviceName, ok = checkAuth(arg["APIKey"]); ok {
		logging.Info("rpc", "["+serviceName+"] Stop("+arg["PushURL"]+")")
	} else {
		logging.Warning("rpc", "Invalid auth for Stop("+arg["PushURL"]+")")
		return jsonrpc2.NewError(403, "Invalid API key")
	}

	if tartmanager.Exists(arg["PushURL"]) {
		err := tartmanager.Stop(arg["PushURL"])
		if err != nil {
			return err
		}
		result.Success = true
		return nil
	}
	return errors.New("Could not find tart")
}

// Init RPC creates a tart's metadata without it actually existing yet.
func (t *Tarts) Init(arg map[string]string, result *ArbitrarySuccessResult) error {
	var serviceName string
	var ok bool
	if serviceName, ok = checkAuth(arg["APIKey"]); ok {
		logging.Info("rpc", "["+serviceName+"] Init("+arg["PushURL"]+", "+arg["User"]+")")
	} else {
		logging.Warning("rpc", "Invalid auth for Init("+arg["PushURL"]+")")
		return jsonrpc2.NewError(403, "Invalid API key")
	}

	if tartmanager.Exists(arg["PushURL"]) {
		return errors.New("Tart exists")
	}
	err := tartmanager.PreGitRecieve(arg["PushURL"], arg["User"])
	if err != nil {
		return err
	}
	tartmanager.New(arg["PushURL"], arg["User"])
	result.Success = true
	return nil
}

// AddOwner RPC adds an owner to an existing tart.
func (t *Tarts) AddOwner(arg map[string]string, result *ArbitrarySuccessResult) error {
	var serviceName string
	var ok bool
	if serviceName, ok = checkAuth(arg["APIKey"]); ok {
		logging.Info("rpc", "["+serviceName+"] AddOwner("+arg["PushURL"]+", "+arg["Username"]+")")
	} else {
		logging.Warning("rpc", "Invalid auth for AddOwner("+arg["PushURL"]+")")
		return jsonrpc2.NewError(403, "Invalid API key")
	}

	if !tartmanager.Exists(arg["PushURL"]) {
		return errors.New("Tart does not exist")
	}
	tart := tartmanager.Get(arg["PushURL"])

	for _, o := range tart.Owners {
		if o == arg["Username"] {
			return errors.New("User is already an owner")
		}
	}

	tart.Owners = append(tart.Owners, arg["Username"])
	tartmanager.Save(arg["PushURL"], tart)
	return nil
}

func setEnv(envList []string, envString, delString string) []string {
	key := strings.Split(envString, "=")[0]
	var output []string

	for _, envEntry := range envList {
		if (strings.Split(envEntry, "=")[0] == key) || strings.Split(envEntry, "=")[0] == delString {
			//no op
		} else {
			output = append(output, envEntry)
		}
	}
	if envString != "" {
		output = append(output, envString)
	}
	return output
}
