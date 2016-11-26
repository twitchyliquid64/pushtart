package tartmanager

import (
	"pushtart/config"
	"pushtart/logging"
	"sync"
	"time"

	gsig "github.com/jondot/gosigar"
)

// PIDs in this map are not checked to see if they are stil running.
var sentryBlacklist = map[int]bool{}
var sentryLock sync.Mutex

func blacklistPidFromSentry(pid int) {
	sentryLock.Lock()
	defer sentryLock.Unlock()
	sentryBlacklist[pid] = true
}

func removePidFromSentryBlacklist(pid int) {
	sentryLock.Lock()
	defer sentryLock.Unlock()
	if _, ok := sentryBlacklist[pid]; ok {
		delete(sentryBlacklist, pid)
	}
}

func pidInSentryBlacklist(pid int) bool {
	_, exists := sentryBlacklist[pid]
	return exists
}

func runSentry() {
	sentryLock.Lock()
	defer sentryLock.Unlock()
	tarts := config.All().Tarts

	for pushURL, tart := range tarts {
		if tart.IsRunning && tart.PID > 0 && !pidInSentryBlacklist(tart.PID) {
			ps := gsig.ProcState{}
			if err := ps.Get(tart.PID); err != nil {
				logging.Warning("run-sentry", "Error getting process info for "+pushURL+": "+err.Error())
				logging.Warning("run-sentry", "Cleaning up execution.")
				sentryLock.Unlock()
				err = Stop(pushURL)
				sentryLock.Lock()
				if err != nil {
					logging.Error("run-sentry", err.Error())
				}
				continue
			}
		}
	}
}

// RunSentry is the entrypoint to the runsentry goroutine.
func RunSentry() {
	for {
		runSentry()
		time.Sleep(time.Second * time.Duration(config.All().RunSentryInterval))
	}
}
