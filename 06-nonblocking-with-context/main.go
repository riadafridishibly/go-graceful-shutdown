package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/riadafridishibly/go-graceful-shutdown/utils"
)

func responsive(ctx context.Context) (string, error) {
	type ret struct {
		value string
		err   error
	}
	ch := make(chan ret)
	go func() {
		v, err := utils.BlockingFunc()
		ch <- ret{v, err}
	}()
	select {
	case <-ctx.Done():
		return "", context.Cause(ctx)
	case v := <-ch:
		return v.value, v.err
	}
}

func main() {
	fmt.Println("PID:", os.Getpid())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancelCause(context.Background())

	utils.SimulateSendSignal(1*time.Second, os.Interrupt)

	go func() {
		got := <-sig
		cancel(fmt.Errorf("signal %s", got))
	}()

	v, err := responsive(ctx)
	fmt.Printf("Value: %q, err: %v\n", v, err)
}
