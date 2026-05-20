package demo

import (
	"github.com/dekaiju/go-skeleton/pkg/log"
)

func helloHandler(args []string) error {
	log.Println(args)
	log.Println("Hello")
	return nil
}
