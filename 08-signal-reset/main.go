package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("PID:", os.Getpid())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	got := <-sigCh

	// Comment out this line, and try Pressing `C-c` again
	signal.Reset(os.Interrupt, syscall.SIGTERM)

	go func() {
		for got := range sigCh {
			fmt.Printf("Received Signal: %s, Sig Num: %d\n", got, got)
		}
	}()

	fmt.Printf("Received Signal: %s, Sig Num: %d\n", got, got)
	for i := 0; i < 5; i++ {
		fmt.Printf("Exiting in %d sec\n", 5-i)
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Exited")
}
