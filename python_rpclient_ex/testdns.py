import jsonrpclib
#jsonrpclib.config.version = 1.0
server = jsonrpclib.Server('http://localhost:8080/rpc')
#server.__send('RPCServer.SysStats', [1])

print server.DNSExtension.List(APIKey="frffhi34NqTef3Bx")
print server.DNSExtension.SetA(APIKey="frffhi34NqTef3Bx",Domain="testdomain",Address="192.168.0.1",TTL=100)
print server.DNSExtension.List(APIKey="frffhi34NqTef3Bx")
print server.DNSExtension.DeleteA(APIKey="frffhi34NqTef3Bx",Domain="testdomain")
print server.DNSExtension.List(APIKey="frffhi34NqTef3Bx")
