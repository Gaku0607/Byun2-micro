package tool

import "log"

var ErrChan chan error = make(chan error)

func PrintFailed() {
	if err := <-ErrChan; err != nil {
		log.Fatalln(err.Error())
	}
}
