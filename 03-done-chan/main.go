package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/riadafridishibly/go-graceful-shutdown/utils"
)

func splitStringDone(s string, done <-chan bool) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for _, v := range strings.Fields(s) {
			select {
			case ch <- v:
				// Just for blocking for 1 sec
				select {
				case <-time.After(1 * time.Second):
				case <-done:
					return
				}
			case <-done:
				return
			}
		}
	}()
	return ch
}

func printer(name string, ch <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range ch {
		fmt.Printf("%s: value = %v\n", name, v)
	}
}

func main() {
	fmt.Println("PID:", os.Getpid())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	utils.SimulateSendSignal(2*time.Second, os.Interrupt)

	done := make(chan bool)
	go func() {
		got := <-sigCh
		fmt.Printf("Received Signal: %s, Sig Num: %d\n", got, got)

		// Close the done channel to signal the `splitStringDone` function that
		// we are no longer interested, we're quiting.
		close(done)
	}()

	ch := splitStringDone("a b c d e f g", done)

	var wg sync.WaitGroup

	wg.Add(2)
	go printer("Printer 1", ch, &wg)
	go printer("Printer 2", ch, &wg)

	wg.Wait()

	fmt.Println("Exited!")

	// Print the goroutine stack trace,
	// to check which goroutines are currently alive
	// debug.SetTraceback("all")
	// panic("show me the stacks")
}
