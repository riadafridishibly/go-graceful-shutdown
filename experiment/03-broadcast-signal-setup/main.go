package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"time"
)

const str1 = "v1 v2 v3 v4 v5 v6"
const str2 = "x1 x2 x3 x4 x5 x6"

func producer(name, s string, ch chan<- string) {
	for _, v := range strings.Fields(s) {
		ch <- v
		time.Sleep(1 * time.Second)
	}
}

func consumer(name string, ch <-chan string) {
	for v := range ch {
		fmt.Printf("%s: got value: %v\n", name, v)
	}
}

func main() {
	log.Println("Process PID:", os.Getpid())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	ch := make(chan string)
	go producer("PROD 1", str1, ch)
	go producer("PROD 2", str2, ch)

	go consumer("CONSUMER 1", ch)
	go consumer("CONSUMER 2", ch)

	got := <-sigCh
	n := int(got.(syscall.Signal))
	log.Println("Received Signal: ", n, got.String())

	// Print the goroutine stack trace
	debug.SetTraceback("all")
	panic("show me the stacks")
}
