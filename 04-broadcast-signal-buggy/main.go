package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"strings"
	"sync"
	"syscall"
	"time"
)

const str1 = "v1 v2 v3 v4 v5 v6"
const str2 = "x1 x2 x3 x4 x5 x6"

func producer(name, s string, ch chan<- string, done <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()

	var sendCh chan<- string = ch
	// What's the bug in this code?
	// - Do we get all the values?
	// - When does the loop exit? What if strings.Fields(s) exhausted?
	for _, v := range strings.Fields(s) {
		select {
		case sendCh <- v:
			sendCh = nil
		case <-done:
			return
		case <-tick.C:
			sendCh = ch
		}
	}
}

func consumer(name string, ch <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range ch {
		fmt.Printf("%s: got value: %v\n", name, v)
	}
}

func main() {
	log.Println("Process PID:", os.Getpid())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	done := make(chan bool)
	ch := make(chan string)

	var pwg, cwg sync.WaitGroup // wait group for producers

	pwg.Add(2)
	go producer("PROD 1", str1, ch, done, &pwg)
	go producer("PROD 2", str2, ch, done, &pwg)

	cwg.Add(2)
	go consumer("CONSUMER 1", ch, &cwg)
	go consumer("CONSUMER 2", ch, &cwg)

	got := <-sigCh
	n := int(got.(syscall.Signal))
	log.Println("Received Signal: ", n, got.String())
	close(done)

	pwg.Wait()
	log.Println("All producers returned")
	close(ch) // All producers returned! Close the chan to signal consumers
	cwg.Wait()
	log.Println("All consumers returned")

	// Print the goroutine stack trace
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 2)
}
