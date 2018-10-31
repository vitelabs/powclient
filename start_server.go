package powClient

import (
	"github.com/gin-gonic/gin"
	"github.com/vitelabs/go-vite/log15"
	"powClient/service"
)

var (
	port = "6007"
)

func StartUp(env string) {
	service.InitUrl(env)
	router := gin.New()
	registerRouterPow(router)
	log15.Info("Server start listen in " + port)
	router.Run(":" + port)
}

func registerRouterPow(engine *gin.Engine) {
	router := engine.Group("/api")
	router.POST("/generate_work", service.WorkDetail)
	router.POST("/validate_work", service.VaildDetail)
	router.POST("/cancel_work", service.CancelDetail)
}
