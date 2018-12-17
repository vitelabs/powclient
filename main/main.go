package main

import (
	"flag"
	"github.com/vitelabs/powclient"
	"github.com/vitelabs/powclient/log15"
	"github.com/vitelabs/powclient/metrics"
)

var log = log15.New("module", "main")

func main() {
	flag.Parse()
	InitLog(DefaultDataDir(), "info")

	metrics.InitMetrics(*metricsEnable, *metricsEnable)
	if *metricsEnable == true {
		powclient.SetUpMetrics(*influxDBUrl)
	}
	switch *mtype {
	case "gpu":
		powclient.StartUpGpu(*serverUrl)
	case "cpu":
		powclient.StartUpCpu(*serverUrl)
	default:
		log.Error("processor type error")
		return
	}
}
