package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// const str1 = "v1 v2 v3 v4 v5 v6"
// const str2 = "x1 x2 x3 x4 x5 x6"

// func splitString(s string) <-chan string {
// 	ch := make(chan string)
// 	go func() {
// 		defer close(ch)
// 		for _, v := range strings.Fields(s) {
// 			ch <- v
// 			time.Sleep(1 * time.Second)
// 		}
// 	}()
// 	return ch
// }

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
	fmt.Println("Process PID:", os.Getpid())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

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

	// Print the goroutine stack trace
	// debug.SetTraceback("all")
	// panic("show me the stacks")
}
