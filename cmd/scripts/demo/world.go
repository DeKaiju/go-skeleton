package demo

import (
	"github.com/dekaiju/go-skeleton/pkg/log"
)

func worldHandler(args []string) error {
	log.Println(args)
	log.Println("world")
	return nil
}
