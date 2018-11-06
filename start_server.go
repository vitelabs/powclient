package powClient

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vitelabs/go-vite/log15"
	"powClient/service/cpu"
	"powClient/service/gpu"
)

var (
	gpu_port = "6007"
	cpu_port = "6008"
)

func StartUpGpu(env string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Print("Gpu panic", err)
		}
	}()
	gpu.InitUrl(env)
	router := gin.New()
	registerRouterPowGpu(router)
	log15.Info("Server start listen in " + gpu_port)
	router.Run(":" + gpu_port)
}

func StartUpCpu(env string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Print("Cpu panic", err)
		}
	}()
	router := gin.New()
	registerRouterPowCpu(router)
	log15.Info("Server start listen in " + cpu_port)
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
