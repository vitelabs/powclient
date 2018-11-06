package main

import (
	"flag"
	"fmt"
	"github.com/vitelabs/powclient"
)

var (
	env   = flag.String("env", "127:0:0:1", "env ip")
	mtype = flag.String("type", "cpu", "machine type")
)

func main() {
	flag.Parse()
	if *mtype == "gpu" {
		powclient.StartUpGpu(*env)
	}
	if *mtype == "cpu" {
		powclient.StartUpCpu(*env)
	}
	fmt.Print("type error")
}
