package main

import (
	"github.com/keypeearr/romanizer/src/server"
)

func main() {
	if err := server.Run(); err != nil {
		panic(err)
	}
}
