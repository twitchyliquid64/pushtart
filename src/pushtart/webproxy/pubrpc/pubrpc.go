package pubrpc

import (
	"errors"
	"pushtart/config"
	"pushtart/dnsserv"

	sigar "github.com/cloudfoundry/gosigar"
)

//SysStatsArgs represents the arguments passed to SysStats RPC.
type SysStatsArgs struct {
	A int
}

//SysStatsResult represents the result of a successful SysStats RPC.
type SysStatsResult struct {
	Name      string
	Load      sigar.LoadAverage
	Uptime    sigar.Uptime
	Mem       sigar.Mem
	Swap      sigar.Swap
	CacheUsed int
}

//RPCService represents the RPC server available via a webproxy URI.
type RPCService int

//SysStats RPC returns a structure representing the state of the machines resources.
func (t *RPCService) SysStats(args int, result *SysStatsResult) error {
	concreteSigar := sigar.ConcreteSigar{}

	result.Mem.Get()
	result.Swap.Get()
	result.Uptime.Get()
	result.Name = config.All().Name
	result.CacheUsed = dnsserv.GetCacheUsed()

	avg, err := concreteSigar.GetLoadAverage()
	if err != nil {
		return errors.New("Failed to get load average")
	}
	result.Load = avg

	return nil
}
