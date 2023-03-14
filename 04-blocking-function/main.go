package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/riadafridishibly/go-graceful-shutdown/utils"
)

func nonresponsive(done <-chan bool) (string, error) {
	select {
	case <-done:
		return "", errors.New("operation cancelled")
	default:
		return utils.BlockingFunc() // select won't do anything
	}
}

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	utils.SimulateSendSignal(1*time.Second, os.Interrupt)

	done := make(chan bool)
	go func() {
		<-sig
		close(done)
	}()

	v, err := nonresponsive(done)
	if err == nil {
		fmt.Println(">>> CANCEL DID NOT WORK")
	}
	fmt.Printf("Value: %q, err: %v\n", v, err)
}
