package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/riadafridishibly/go-graceful-shutdown/utils"
)

func main() {
	fmt.Println("PID:", os.Getpid())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	utils.SimulateSendSignal(1*time.Second, os.Interrupt)
	utils.SimulateSendSignal(2*time.Second, syscall.SIGTERM)
	utils.SimulateSendSignal(3*time.Second, os.Interrupt)

	got := <-sigCh
	fmt.Printf("Received Signal: %s, Sig Num: %d\n", got, got)

	// Uncomment the next line and run the program again
	// signal.Reset(os.Interrupt, syscall.SIGTERM)

	go func() {
		// To show that we're still receiving signals
		for got := range sigCh {
			fmt.Printf("Received Signal: %s, Sig Num: %d\n", got, got)
		}
	}()

	for i := 0; i < 5; i++ {
		fmt.Printf("Exiting in %d sec\n", 5-i)
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Exited")
}
