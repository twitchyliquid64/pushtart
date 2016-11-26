package tartmanager

// Run manager gathers statistics on running tarts, storing them in memory.
// Sentry: makes sure that running processes whose PID corresponds to a tart execution
// are still running, updatingb pushtart state otherwise.
// Metrics about running tarts.
var metrics = map[string]RunMetrics{}

// RunMetrics stores process level information about a running tart.
type RunMetrics struct {
	PID     int
	CmdLine []string
}
