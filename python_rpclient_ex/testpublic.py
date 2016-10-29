import jsonrpclib
#jsonrpclib.config.version = 1.0
server = jsonrpclib.Server('http://localhost:8080/pubrpc')
#server.__send('RPCServer.SysStats', [1])
print server.RPCService.SysStats()['Uptime']['Length']
