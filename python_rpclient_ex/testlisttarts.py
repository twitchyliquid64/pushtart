import jsonrpclib
#jsonrpclib.config.version = 1.0
server = jsonrpclib.Server('http://localhost:8080/rpc')
#server.__send('RPCServer.SysStats', [1])
print server.Service.ListTarts(APIKey="frffhi34NqTef3Bx")
