package rpc

import (
	"bytes"
	"errors"
	"pushtart/config"
	"pushtart/logging"
	"pushtart/sshserv/cmd_registry"
	"reflect"
	"strconv"
	"strings"

	"github.com/powerman/rpc-codec/jsonrpc2"
)

//AuthenticationArgument is the default argument for proceedures which require no arguments but authentication.
type AuthenticationArgument struct {
	APIKey string
}

//Service represents the authenticated RPC server for general methods, available via a webproxy URI.
type Service int

//ListTartsResult represents the result of a successful ListTarts RPC.
type ListTartsResult struct {
	Tarts map[string]config.Tart
}

//ListTarts RPC returns a list of tarts.
func (t *Service) ListTarts(arg *AuthenticationArgument, result *ListTartsResult) error {
	var serviceName string
	var ok bool
	if serviceName, ok = checkAuth(arg.APIKey); ok {
		logging.Info("rpc", "["+serviceName+"] ListTarts()")
	} else {
		logging.Warning("rpc", "Invalid auth for ListTarts()")
		return jsonrpc2.NewError(403, "Invalid API key")
	}

	result.Tarts = map[string]config.Tart{}
	for name, t := range config.All().Tarts {
		result.Tarts[name] = t
	}
	return nil
}

//ListUsersResult represents the result of a successful ListTarts RPC.
type ListUsersResult struct {
	Users map[string]config.User
}

//ListUsers RPC returns a list of users.
func (t *Service) ListUsers(arg *AuthenticationArgument, result *ListUsersResult) error {
	var serviceName string
	var ok bool
	if serviceName, ok = checkAuth(arg.APIKey); ok {
		logging.Info("rpc", "["+serviceName+"] ListUsers()")
	} else {
		logging.Warning("rpc", "Invalid auth for ListUsers()")
		return jsonrpc2.NewError(403, "Invalid API key")
	}

	result.Users = map[string]config.User{}
	for name, t := range config.All().Users {
		temp := t
		temp.Password = "<nil>"
		result.Users[name] = temp
	}
	return nil
}

//RunCommandArgument represents the parameters passed to the RunCommand RPC.
type RunCommandArgument struct {
	APIKey  string
	Command string
	Args    map[string]string
}

//RunCommandResult represents the datastructure returned upon successful completion of the RunCommand RPC.
type RunCommandResult struct {
	Output string
}

//RunCommand RPC executes the given command and returns output.
func (t *Service) RunCommand(arg *RunCommandArgument, result *RunCommandResult) error {
	var serviceName string
	var ok bool
	if serviceName, ok = checkAuth(arg.APIKey); ok {
		logging.Info("rpc", "["+serviceName+"] RunCommand("+arg.Command+")")
	} else {
		logging.Warning("rpc", "Invalid auth for RunCommand("+arg.Command+")")
		return jsonrpc2.NewError(403, "Invalid API key")
	}

	if ok, cmd := cmd_registry.Command(arg.Command); ok {
		buf := new(bytes.Buffer)
		cmd(arg.Args, buf, "")
		result.Output = buf.String()
	} else {
		return errors.New("Command not found")
	}

	return nil
}

//GetConfigValueArgument represents the parameters passed to the GetConfigValue RPC.
type GetConfigValueArgument struct {
	APIKey string
	Field  string
}

//GetConfigValueResult represents the result of a successful GetConfigValue RPC.
type GetConfigValueResult struct {
	Value   interface{}
	Error   string
	Success bool
}

//GetConfigValue RPC returns the structure at the given field in the configuration.
func (t *Service) GetConfigValue(arg *GetConfigValueArgument, result *GetConfigValueResult) error {
	var serviceName string
	var ok bool
	if serviceName, ok = checkAuth(arg.APIKey); ok {
		logging.Info("rpc", "["+serviceName+"] GetConfigValue("+arg.Field+")")
	} else {
		logging.Warning("rpc", "Invalid auth for GetConfigValue("+arg.Field+")")
		return jsonrpc2.NewError(403, "Invalid API key")
	}

	var err error
	result.Value, err = getVal(arg.Field, reflect.ValueOf(config.All()).Elem())
	if err != nil {
		result.Error = err.Error()
		result.Success = false
	} else {
		result.Success = true
	}
	return nil
}

//SetConfigValueArgument represents the parameters passed to the SetConfigValue RPC.
type SetConfigValueArgument struct {
	APIKey string
	Field  string
	Value  string
}

//SetConfigValueResult represents the result of a successful SetConfigValue RPC.
type SetConfigValueResult struct {
	Error   string
	Success bool
}

//SetConfigValue RPC returns the structure at the given field in the configuration.
func (t *Service) SetConfigValue(arg *SetConfigValueArgument, result *SetConfigValueResult) error {
	var serviceName string
	var ok bool
	if serviceName, ok = checkAuth(arg.APIKey); ok {
		logging.Info("rpc", "["+serviceName+"] SetConfigValue("+arg.Field+")")
	} else {
		logging.Warning("rpc", "Invalid auth for SetConfigValue("+arg.Field+")")
		return jsonrpc2.NewError(403, "Invalid API key")
	}

	err := setVal(arg.Field, arg.Value, reflect.ValueOf(config.All()).Elem())
	if err != nil {
		result.Error = err.Error()
		result.Success = false
	} else {
		result.Success = true
	}
	return nil
}

func getVal(query string, val reflect.Value) (out interface{}, err error) {
	spl := strings.Split(query, ".")
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		valueField := val.Field(i)
		tag := typeField.Tag
		if len(spl) == 1 { //no other sections like DNS.Listener
			if typeField.Name == spl[0] && tag.Get("getConfigValue") != "block" {
				return valueField.Interface(), nil
			}
		} else {
			if typeField.Name == spl[0] {
				return getVal(strings.Join(spl[1:], "."), valueField)
			}
		}
	}

	return "", errors.New("Could not find field")
}

func setVal(query, newVal string, val reflect.Value) (err error) {
	spl := strings.Split(query, ".")
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		valueField := val.Field(i)
		tag := typeField.Tag
		if len(spl) == 1 { //no other sections like DNS.Listener
			if typeField.Name == spl[0] && tag.Get("getConfigValue") != "block" {
				v, err := strToVal(newVal, valueField)
				if err != nil {
					return err
				}
				if !valueField.CanSet() {
					return errors.New("Cannot set field: " + typeField.Name)
				}
				valueField.Set(v)
				return nil
			}
		} else {
			if typeField.Name == spl[0] {
				return setVal(strings.Join(spl[1:], "."), newVal, valueField)
			}
		}
	}

	return errors.New("Could not find field")
}

func strToVal(in string, template reflect.Value) (reflect.Value, error) {
	switch template.Kind() {
	case reflect.Bool:
		vb, err := strconv.ParseBool(in)
		return reflect.ValueOf(vb), err
	case reflect.String:
		return reflect.ValueOf(in), nil
	case reflect.Int64:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Int:
		vi, err := strconv.Atoi(in)
		return reflect.ValueOf(vi), err
	}
	return reflect.ValueOf(nil), errors.New("Don't know how to process " + template.Kind().String())
}

func checkAuth(key string) (service string, ok bool) {
	for _, entry := range config.All().APIKeys {
		if entry.Key == key {
			return entry.Service, true
		}
	}
	return "", false
}
