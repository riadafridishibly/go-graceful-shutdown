package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"sync"
	"syscall"
	"time"
)

const str1 = "v1 v2 v3 v4 v5 v6"
const str2 = "x1 x2 x3 x4 x5 x6"

func producer(ctx context.Context, name, s string, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()

	items := strings.Fields(s)

	var sendCh chan<- string = ch

	for len(items) > 0 {
		v := items[0]
		select {
		case sendCh <- v:
			sendCh = nil // block sending
			items = items[1:]
		case <-ctx.Done():
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ch := make(chan string)

	var pwg, cwg sync.WaitGroup // wait group for producers

	pwg.Add(2)
	go producer(ctx, "PRODUCER 1", str1, ch, &pwg)
	go producer(ctx, "PRODUCER 2", str2, ch, &pwg)

	cwg.Add(2)
	go consumer("CONSUMER 1", ch, &cwg)
	go consumer("CONSUMER 2", ch, &cwg)

	<-ctx.Done()
	stop()
	log.Println("Received cancellation signal")
	log.Println("But we don't know which signal triggered the cancellation")
	log.Println("Because the signal was thrown away")
	log.Println("See here:", "https://cs.opensource.google/go/go/+/refs/tags/go1.20.2:src/os/signal/signal.go;l=289")

	pwg.Wait()
	log.Println("All producers returned")
	close(ch) // All producers returned! Close the chan to signal consumers
	cwg.Wait()
	log.Println("All consumers returned")

	// Print the goroutine stack trace
	debug.SetTraceback("all")
	panic("show me the stacks")
}
