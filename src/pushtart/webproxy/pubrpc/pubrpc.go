package pubrpc

import (
	"errors"
	"os"
	"pushtart/config"
	"pushtart/dnsserv"
	"syscall"

	sigar "github.com/cloudfoundry/gosigar"
)

func getDiskFreeAndTotal() (uint64, uint64) {
	var stat syscall.Statfs_t
	wd, _ := os.Getwd()
	syscall.Statfs(wd, &stat)
	// Available blocks * size per block = available space in bytes
	return stat.Bavail * uint64(stat.Bsize), stat.Blocks * uint64(stat.Bsize)
}

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
	DiskFree  uint64
	DiskTotal uint64
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
	result.DiskFree, result.DiskTotal = getDiskFreeAndTotal()

	avg, err := concreteSigar.GetLoadAverage()
	if err != nil {
		return errors.New("Failed to get load average")
	}
	result.Load = avg

	return nil
}
