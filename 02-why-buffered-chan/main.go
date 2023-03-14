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
	sigCh := make(chan os.Signal, 1) // Change this to unbuffered, make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	utils.SimulateSendSignal(1*time.Second, os.Interrupt)

	utils.BlockingFunc()

	got := <-sigCh
	fmt.Printf("Received Signal: %s, Sig Num: %d\n", got, got)
}
