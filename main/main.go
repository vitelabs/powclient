package main

import (
	"flag"
	"fmt"
	"powClient"
)

var (
	env   = flag.String("env", "127:0:0:1", "env ip")
	mtype = flag.String("type", "gpu", "machine type")
)

func main() {
	flag.Parse()
	if *mtype == "gpu" {
		powClient.StartUpGpu(*env)
	}
	if *mtype == "cpu" {
		powClient.StartUpCpu(*env)
	}
	fmt.Print("type error")
}
