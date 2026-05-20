package main

import (
	"runtime"
	"strconv"

	"github.com/dekaiju/go-skeleton/cmd"
	"github.com/dekaiju/go-skeleton/pkg/log"
)

var (
	SetCpuCount string
)

// @title go-skeleton
// @version 2.0
// @description Go service skeleton with a reusable backend project structure
// @termsOfService https://github.com/dekaiju/go-skeleton
func main() {
	if SetCpuCount != "" {
		procsNum, err := strconv.Atoi(SetCpuCount)
		if err == nil {
			runtime.GOMAXPROCS(procsNum)
			log.Printf("GOMAXPROCS set to %d", procsNum)
		}
	}
	cmd.Execute()
}
