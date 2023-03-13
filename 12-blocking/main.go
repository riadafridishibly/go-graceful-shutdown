package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func blockingFunc() (string, error) {
	fmt.Println("Blocking func started, will sleep for 5 sec")
	defer fmt.Println("Blocking func finished")

	time.Sleep(10 * time.Second)
	return "some value", nil
}

func nonresponsive(done <-chan bool) (string, error) {
	select {
	case <-done:
		return "", errors.New("operation cancelled")
	default:
		return blockingFunc() // select won't do anything
	}
}

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	done := make(chan bool)
	go func() {
		<-sig
		close(done)
	}()
	v, err := nonresponsive(done)
	fmt.Printf("Value: %q, err: %v\n", v, err)
}
