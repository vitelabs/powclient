package main

import (
	"flag"
	"powClient"
)

var env = flag.String("env", "127:0:0:1", "env ip")

func main() {
	flag.Parse()
	powClient.StartUp(*env)
}
