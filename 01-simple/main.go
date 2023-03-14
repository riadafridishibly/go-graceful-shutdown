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

	// Wait for signal
	got := <-sigCh

	fmt.Printf("Received Signal: %s, Sig Num: %d\n", got, got)
}
