// Go port of Coda Hale's Metrics library
//
// <https://github.com/rcrowley/go-metrics>
//
// Coda Hale's original work: <https://github.com/codahale/metrics>
package metrics

import (
	"github.com/elastic/gosigar"
	"github.com/vitelabs/powclient/log15"
	"runtime"
	"strings"
	"time"
)

// MetricsEnabled is checked by the constructor functions for all of the
// standard metrics.  If it is true, the metric returned is a stub.
//
// This global kill-switch helps quantify the observer effect and makes
// for less cluttered pprof profiles.
var (
	MetricsEnabled     = false
	InfluxDBExportFlag = false
	log                = log15.New("module", "metrics")
)

func InitMetrics(metricFlag, influxDBFlag bool) {
	MetricsEnabled = metricFlag
	InfluxDBExportFlag = influxDBFlag
}

var (
	CodexecRegistry       = NewPrefixedChildRegistry(DefaultRegistry, "/codexec")
	TimeConsumingRegistry = NewPrefixedChildRegistry(CodexecRegistry, "/timeconsuming")
)

func TimeConsuming(names []string, sinceTime time.Time) {
	if !MetricsEnabled {
		return
	}
	var name string
	for _, v := range names {
		name += "/" + strings.ToLower(v)
	}
	if timer, ok := GetOrRegisterResettingTimer(name, TimeConsumingRegistry).(*StandardResettingTimer); timer != nil && ok {
		timer.UpdateSince(sinceTime)
	}
}

var (
	systemRegistry = NewPrefixedChildRegistry(DefaultRegistry, "/system")
	cpuRegistry    = NewPrefixedChildRegistry(systemRegistry, "/cpu")
)

func RefreshCpuStats(refresh time.Duration, prevProcessCPUTime float64, prevSystemCPUUsage gosigar.Cpu) (float64, gosigar.Cpu) {
	if !MetricsEnabled {
		return 0, gosigar.Cpu{}
	}
	frequency := float64(refresh / time.Second)
	numCPU := float64(runtime.NumCPU())
	curSystemCPUUsage := gosigar.Cpu{}
	var curProcessCPUTime float64

	if processCPUTimeGuage, ok := GetOrRegisterGaugeFloat64("/processtime", cpuRegistry).(*StandardGaugeFloat64); ok && processCPUTimeGuage != nil {
		curProcessCPUTime = getProcessCPUTime()
		deltaProcessCPUTime := curProcessCPUTime - prevProcessCPUTime
		processCPUTime := deltaProcessCPUTime / frequency / numCPU * 100

		processCPUTimeGuage.Update(processCPUTime)
	}
	if systemCPUUsageGuage, ok := GetOrRegisterGaugeFloat64("/sysusage", cpuRegistry).(*StandardGaugeFloat64); ok && systemCPUUsageGuage != nil {
		curSystemCPUUsage.Get()
		deltaSystemCPUUsage := curSystemCPUUsage.Delta(prevSystemCPUUsage)
		systemCPUValue := float64(deltaSystemCPUUsage.Sys+deltaSystemCPUUsage.User) / frequency / numCPU

		systemCPUUsageGuage.Update(systemCPUValue)
	}
	return curProcessCPUTime, curSystemCPUUsage
}
