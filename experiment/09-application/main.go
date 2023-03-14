package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func heavyTask() (string, error) {
	time.Sleep(1 * time.Second)
	if v := rand.Intn(2); v == 0 {
		return "", errors.New("value is 0")
	} else {
		return fmt.Sprintf("value: %d", rand.Intn(100)), nil
	}
}

type App struct {
	closing chan chan error
}

func NewApp() *App {
	return &App{
		closing: make(chan chan error, 1),
	}
}

func (a *App) Run() {
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()

	type taskResult struct {
		res string
		err error
	}

	taskResChan := make(chan taskResult)
	var err error
	for {
		select {
		case errc := <-a.closing:
			errc <- err
			return
		case task := <-taskResChan:
			err = task.err
			fmt.Println("Task: ", task.res, "Error: ", task.err)
		case <-tick.C:
			go func() {
				res, err := heavyTask()
				taskResChan <- taskResult{res, err}
			}()
		}
	}
}

func (a *App) Close() error {
	errc := make(chan error)
	a.closing <- errc
	return <-errc
}

func main() {
	log.Println("Process PID:", os.Getpid())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	app := NewApp()
	go app.Run()

	got := <-sigCh
	n := int(got.(syscall.Signal))
	signal.Reset(os.Interrupt, syscall.SIGTERM)
	log.Println("Received Signal:", n, got.String())
	log.Println("Closed app. err:", app.Close())
}
