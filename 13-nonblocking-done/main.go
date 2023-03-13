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

	time.Sleep(5 * time.Second)
	return "some value", nil
}

func responsive(done <-chan bool) (string, error) {
	type result struct {
		value string
		err   error
	}
	ch := make(chan result)
	go func() {
		v, err := blockingFunc()
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

	done := make(chan bool)
	go func() {
		<-sig
		close(done)
	}()
	v, err := responsive(done)
	fmt.Printf("Value: %q, err: %v\n", v, err)
}
