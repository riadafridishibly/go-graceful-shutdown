package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Process PID:", os.Getpid())
	sigCh := make(chan os.Signal, 1) // Change this to unbuffered
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	fmt.Println("Sleep started. Now press C-c")
	time.Sleep(10 * time.Second)
	fmt.Println("Sleep done...")

	got := <-sigCh
	fmt.Printf("Received Signal: %s, Sig Num: %d\n", got, got)
}
