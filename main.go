package main

import (
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"github.com/dsbarabash/shopping-lists/internal/service"
	"sync"
	"time"
)

func main() {
	ch := make(chan interface{})
	stop := time.NewTimer(10 * time.Second)
	defer close(ch)
	logging := time.NewTicker(time.Millisecond * 200)
	tick := time.NewTicker(time.Millisecond * 150)
	defer tick.Stop() // освободим ресурсы, при завершении работы функции
	defer logging.Stop()
	for {
		select {
		case <-stop.C:
			// stop - Timer, который через 10 секунд даст сигнал завершить работу

			return
		case <-tick.C:
			// tick - Ticker, посылающий сигнал выполнить работу каждые 150 миллисекунд
			wg := new(sync.WaitGroup)
			wg.Add(2)
			go service.CreateRandomStructs(ch)

			go repository.CheckInterface2(ch)
			wg.Done()
		case <-logging.C:
			// tick - Ticker, посылающий сигнал выполнить лоигрвоание каждые 200 миллисекунд
			go repository.LoggingSlice()
		}
	}
}
