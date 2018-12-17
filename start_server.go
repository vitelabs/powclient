package powclient

import (
	"github.com/gin-gonic/gin"
	"github.com/vitelabs/powclient/log15"
	"github.com/vitelabs/powclient/metrics"
	"github.com/vitelabs/powclient/metrics/influxdb"
	"github.com/vitelabs/powclient/service/cpu"
	"github.com/vitelabs/powclient/service/gpu"
	"time"
)

var (
	gpu_port = "6007"
	cpu_port = "6008"
	log      = log15.New("module", "startup")
)

func StartUpGpu(env string) {
	defer func() {
		if err := recover(); err != nil {
			log.Info("Gpu panic", nil, err)
		}
	}()
	gpu.InitUrl(env)
	router := gin.New()
	registerRouterPowGpu(router)
	log.Info("Server start listen in " + gpu_port)
	router.Run(":" + gpu_port)
}

func StartUpCpu(env string) {
	defer func() {
		if err := recover(); err != nil {
			log.Info("Cpu panic", nil, err)
		}
	}()
	router := gin.New()
	registerRouterPowCpu(router)
	log.Info("Server start listen in " + cpu_port)
	router.Run(":" + cpu_port)
}

func registerRouterPowGpu(engine *gin.Engine) {
	router := engine.Group("/api")
	router.POST("/generate_work", gpu.WorkDetail)
	router.POST("/cancel_work", gpu.CancelDetail)
	//router.POST("/validate_work", gpu.VaildDetail)
}

func registerRouterPowCpu(engine *gin.Engine) {
	router := engine.Group("/api")
	router.POST("/generate_work", cpu.WorkDetail)
	//router.POST("/validate_work", cpu.VaildDetail)
	//router.POST("/cancel_work", service.CancelDetail)
}

func SetUpMetrics(url string) {
	metrics.InitMetrics(true, true)
	var (
		influxDbDuration  = 10 * time.Second
		influxDbEndpoint  = "http://" + url
		influxDbDatabase  = "metrics"
		influxDbUsername  = "test"
		influxDbPassword  = "test"
		influxDbNamespace = "monitor"
	)

	if metrics.InfluxDBExportFlag == true {
		log.Info("Enabling metrics export to InfluxDB\n")
		go influxdb.InfluxDBWithTags(metrics.DefaultRegistry, influxDbDuration,
			influxDbEndpoint, influxDbDatabase,
			influxDbUsername, influxDbPassword,
			influxDbNamespace, map[string]string{"host": "localhost_pow"})

	} else {
		log.Info("InfluxDBExport disable\n")
	}
}
