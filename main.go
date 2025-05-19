package main

import (
	"fmt"
	"net/http"
	"nova/app"
	"runtime"
)

func init() {
	// start multi-cpu
	core := runtime.NumCPU()
	runtime.GOMAXPROCS(core)
	// start debug pprof
	go func() {
		_ = http.ListenAndServe(":10080", nil)
	}()
}

func main() {
	fmt.Println("The Nova Project")
	nova := app.New()
	nova.Init()
	nova.Start()
}
