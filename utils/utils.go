package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

func SimulateSendSignal(after time.Duration, sig os.Signal) {
	go func() {
		pid := os.Getpid()
		p, err := os.FindProcess(pid)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(after)
		fmt.Printf("==== Sending signal %q to PID(%d)\n", sig, pid)
		if err := p.Signal(sig); err != nil {
			log.Fatal(err)
		}
	}()
}

func BlockingFunc() (string, error) {
	n := 5 * time.Second
	fmt.Printf("Blocking func started, will sleep for %v\n", n)
	defer fmt.Println("Blocking func finished")

	time.Sleep(n)
	return "foo bar baz", nil
}
