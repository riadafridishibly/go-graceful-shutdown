package main

import (
	"context"
	"log"
	"os"
	"os/exec"
)

func Execute() {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "../cmd/longprocess/longprocess")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Fatal("Err starting process")
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatal("Err waiting")
	}
}

func main() {
	Execute()
}
