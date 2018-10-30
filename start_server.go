package powClient

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vitelabs/go-vite/log15"
	"os/exec"
	"powClient/service"
)

var (
	port   = "6007"
	powLog = log15.New("module", "pow_client")
)

func StartUp(env string) {
	service.InitUrl(env)
	router := gin.New()
	registerRouterPow(router)
	log15.Info("Server start listen in " + port)
	router.Run(":" + port)
	runComm()
}

func registerRouterPow(engine *gin.Engine) {
	router := engine.Group("/api")
	router.POST("/generate_work", service.WorkDetail)
	router.POST("/validate_work", service.VaildDetail)
	router.POST("/cancel_work", service.CancelDetail)
}

func runComm() {
	cpu_thread := "4 "
	listen_address := "0.0.0.0:7076"
	cmd1 := exec.Command("cd /Users/crzn/workplace/PycharmProject/dpow-work-server/target/release", "")
	out1, err := cmd1.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out1))

	cmd2 := exec.Command("nohup", " ./dpow-work-server --cpu-threads ", cpu_thread, " --listen-address ", listen_address,
		" > log/out.file 2>&1 &")
	out, err := cmd2.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}
