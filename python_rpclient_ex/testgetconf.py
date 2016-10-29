import jsonrpclib
#jsonrpclib.config.version = 1.0
server = jsonrpclib.Server('http://localhost:8080/rpc')
#server.__send('RPCServer.SysStats', [1])
v = server.Service.GetConfigValue(APIKey="frffhi34NqTef3Bx",Field="DNS.Listener")
print v
print server.Service.SetConfigValue(APIKey="frffhi34NqTef3Bx",Field="DNS.Listener",Value=":5050")
print server.Service.GetConfigValue(APIKey="frffhi34NqTef3Bx",Field="DNS.Listener")
print server.Service.SetConfigValue(APIKey="frffhi34NqTef3Bx",Field="DNS.Listener",Value=v['Value'])
print server.Service.GetConfigValue(APIKey="frffhi34NqTef3Bx",Field="DNS.Listener")
