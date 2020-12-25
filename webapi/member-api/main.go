package main

import (
	"log"

	"github.com/Gaku0607/Byun2-micro/webapi/member-api/runner"
	. "github.com/Gaku0607/Byun2-micro/webapi/tool"
)

func main() {

	if err := runner.Init(); err != nil {
		log.Fatal(err.Error())
	}

	go runner.Run()

	PrintFailed()
}
