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

func responsive(done <-chan bool) (string, error) {
	type result struct {
		value string
		err   error
	}
	ch := make(chan result)
	go func() {
		v, err := utils.BlockingFunc()
		ch <- result{v, err}
	}()
	select {
	case <-done:
		return "", errors.New("process cancelled")
	case v := <-ch:
		return v.value, v.err
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
	v, err := responsive(done)
	fmt.Printf("Value: %q, err: %v\n", v, err)
}
