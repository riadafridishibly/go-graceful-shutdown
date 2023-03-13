package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	errlog = log.New(os.Stderr, "ERR: ", log.LstdFlags)
	outlog = log.New(os.Stdout, "INF: ", log.LstdFlags)
)

func main() {
	outlog.Println("Program started")
	defer outlog.Println("Program exited")
	exit := make(chan bool)
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		close(exit)
	}()
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()
	for {
		select {
		case <-exit:
			return
		case <-tick.C:
			n := rand.Intn(3)
			if n == 0 {
				errlog.Println("Bad value. value =", n)
			} else {
				outlog.Println("Value is", n)
			}
		}
	}
}
