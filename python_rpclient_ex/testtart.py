import jsonrpclib
#jsonrpclib.config.version = 1.0
server = jsonrpclib.Server('http://localhost:8080/rpc')
#server.__send('RPCServer.SysStats', [1])

print server.Tarts.Init(APIKey="frffhi34NqTef3Bx",PushURL="/testRPC",User="xxx")
print server.Tarts.GetTart(APIKey="frffhi34NqTef3Bx",PushURL="/testRPC")
print server.Tarts.EnableOutputLogging(APIKey="frffhi34NqTef3Bx",PushURL="/testRPC", Enable="True")
print server.Tarts.SetName(APIKey="frffhi34NqTef3Bx",PushURL="/testRPC", Name="Test LOL")
print server.Tarts.SetEnv(APIKey="frffhi34NqTef3Bx",PushURL="/testRPC", Key="testkey", Value="TestValue2")
print server.Tarts.GetTart(APIKey="frffhi34NqTef3Bx",PushURL="/testRPC")
print server.Tarts.DelEnv(APIKey="frffhi34NqTef3Bx",PushURL="/testRPC", Key="testkey")
print server.Tarts.GetTart(APIKey="frffhi34NqTef3Bx",PushURL="/testRPC")
print server.Tarts.Start(APIKey="frffhi34NqTef3Bx",PushURL="/testRPC")
print server.Tarts.Stop(APIKey="frffhi34NqTef3Bx",PushURL="/testRPC")
print server.Tarts.Stop(APIKey="frffhi34NqTef3Bx",PushURL="/testRPC")
