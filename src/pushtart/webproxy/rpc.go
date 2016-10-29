package webproxy

import (
	"net/http"
	"net/rpc"
	"pushtart/logging"
	"pushtart/webproxy/pubrpc"
	privrpc "pushtart/webproxy/rpc"

	"github.com/powerman/rpc-codec/jsonrpc2"
)

func pubRPCInit() http.Handler {
	rServ := rpc.NewServer()
	rpcServ := new(pubrpc.RPCService)
	err := rServ.Register(rpcServ)
	if err != nil {
		logging.Error("jsonrpc-init", "rpc.Register() (pub) error: "+err.Error())
	}
	return jsonrpc2.HTTPHandler(rServ)
}

func privRPCInit() http.Handler {
	rServ := rpc.NewServer()
	rpcServ := new(privrpc.Service)
	err := rServ.Register(rpcServ)
	if err != nil {
		logging.Error("jsonrpc-init", "rpc.Register() (priv) error: "+err.Error())
	}
	return jsonrpc2.HTTPHandler(rServ)
}
