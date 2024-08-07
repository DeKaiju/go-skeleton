package main

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/dekaiju/go-skeleton/cmd"
)

var (
	SetCpuCount string
)

// @title go-skeleton
// @version 1.0
// @description Golang skeleton
// @termsOfService https://github.com/dekaiju/go-skeleton
func main() {
	if SetCpuCount != "" {
		procsNum, err := strconv.Atoi(SetCpuCount)
		if err == nil {
			runtime.GOMAXPROCS(procsNum)
			fmt.Printf("GOMAXPROCS num set %v\n", procsNum)
		}
	}
	cmd.Execute()
}
