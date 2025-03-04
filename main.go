package main

import (
	"context"
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"github.com/dsbarabash/shopping-lists/internal/service"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	ch := make(chan interface{})

	defer close(ch)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	defer stop()

	wg := new(sync.WaitGroup)
	wg.Add(3)

	go func() {
		generator(ctx, ch, time.Millisecond*100)
		wg.Done()
	}()
	go func() {
		asyncStore(ctx, ch)
		wg.Done()
	}()
	go func() {
		asyncLogger(ctx, time.Millisecond*200)
		wg.Done()
	}()
	wg.Wait()

}

func generator(ctx context.Context, ch chan any, d time.Duration) {
	ticker := time.NewTicker(d)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutdown")
			return
		case <-ticker.C:
			ch <- service.CreateRandomStructs()
		}
	}
}

func asyncStore(ctx context.Context, ch chan any) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutdown")
			return
		case data := <-ch:
			repository.CheckInterface(data)
		}
	}
}

func asyncLogger(ctx context.Context, d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutdown")
			return
		case <-ticker.C:
			repository.LoggingSlice()
		}
	}
}
