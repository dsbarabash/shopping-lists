package main

import (
	"context"
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"github.com/dsbarabash/shopping-lists/internal/service"
	"sync"
	"time"
)

func main() {
	ch := make(chan interface{})

	defer close(ch)
	ctx := context.Background()
	ctxWithTimeout, cancelFunction := context.WithTimeout(ctx, time.Duration(1550)*time.Millisecond)

	defer func() {
		cancelFunction()
	}()

	wg := new(sync.WaitGroup)
	wg.Add(3)

	go func() {
		generator(ctxWithTimeout, ch, time.Millisecond*100)
		wg.Done()
	}()
	go func() {
		asyncStore(ctxWithTimeout, ch)
		wg.Done()
	}()
	go func() {
		asyncLogger(ctxWithTimeout, time.Millisecond*200)
		wg.Done()
	}()
	wg.Wait()

}

func generator(ctx context.Context, ch chan any, d time.Duration) {
	ticker := time.NewTicker(d)
	for {
		select {
		case <-ctx.Done():
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
			return
		case <-ticker.C:
			repository.LoggingSlice()
		}
	}
}
