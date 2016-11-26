package tartmanager

import (
	"pushtart/logging"

	gsig "github.com/jondot/gosigar"
)

// RunMetrics stores process level information about a running tart.
type RunMetrics struct {
	PID      int
	State    gsig.ProcState
	Mem      gsig.ProcMem
	Time     gsig.ProcTime
	Children []*RunMetrics
}

// GetStats returns a RunMetrics struct, which describe the running state of a tart.
func GetStats(pushURL string) (*RunMetrics, error) {
	if !Exists(pushURL) {
		return nil, ErrTartNotFound
	}
	tart := Get(pushURL)
	if !tart.IsRunning {
		return nil, ErrTartWrongState
	}
	return getStats(tart.PID)
}

func getStats(pid int) (*RunMetrics, error) {
	ret := &RunMetrics{PID: pid}

	if err := ret.State.Get(pid); err != nil {
		return nil, err
	}
	if err := ret.Mem.Get(pid); err != nil {
		return nil, err
	}
	if err := ret.Time.Get(pid); err != nil {
		return nil, err
	}

	c, err := getChildStats(pid)
	if err != nil {
		return nil, err
	}
	ret.Children = c

	ret.sumReduce(c)
	return ret, nil
}

func getChildStats(pid int) ([]*RunMetrics, error) {
	var children []*RunMetrics
	procs := gsig.ProcList{}
	err := procs.Get()
	if err != nil {
		return nil, err
	}

	for _, cPid := range procs.List {
		processInfo := gsig.ProcState{}
		if err := processInfo.Get(cPid); err != nil {
			logging.Warning("tartmanager-stats", "ProcState.Get(): ", err)
			continue
		}
		if processInfo.Ppid == pid {
			c, err2 := getStats(cPid)
			if err2 != nil {
				return nil, err2
			}
			children = append(children, c)
		}
	}
	return children, nil
}

func (m *RunMetrics) sumReduce(children []*RunMetrics) {
	for _, child := range children {
		m.Mem.MajorFaults += child.Mem.MinorFaults
		m.Mem.MinorFaults += child.Mem.MinorFaults
		m.Mem.PageFaults += child.Mem.PageFaults
		m.Mem.Resident += child.Mem.Resident
		m.Mem.Share += child.Mem.Share
		m.Mem.Size += child.Mem.Size

		m.Time.Sys += child.Time.Sys
		m.Time.Total += child.Time.Total
		m.Time.User += child.Time.User
	}
}
