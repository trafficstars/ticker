package ticker

// This ticker supports changing the interval while running

import (
	"time"
)

type ticker struct {
	ticker *time.Ticker
	C chan time.Time
	stopSignal chan struct{}
	lastTriggerTime time.Time
	runned bool
}

func New(interval time.Duration) *ticker {
	result := &ticker{}
	result.C = make(chan time.Time)
	result.stopSignal = make(chan struct{})
	result.start(interval)
	return result
}

func (ticker *ticker) start(interval time.Duration) {
	ticker.lastTriggerTime = time.Now()
	ticker.ticker = time.NewTicker(interval)
	ticker.run()
}

func (ticker *ticker) run() {
	go ticker.loop()
}

func (ticker *ticker) loop() {
	for {
		select {
		case <-ticker.stopSignal:
			return
		case t := <-ticker.ticker.C:
			ticker.C <- t
		}
	}
}

func (ticker *ticker) flushChan() {
	for {
		select {
		case <-ticker.C:
			continue
		default:
		}
		break
	}
}

func (ticker *ticker) stop() {
	ticker.flushChan()
	ticker.stopSignal <- struct{}{}
	ticker.ticker.Stop()
}

func (ticker *ticker) Stop() {
	ticker.stop()
}

func (ticker *ticker) Restart(newInterval time.Duration) {
	ticker.stop()
	ticker.start(newInterval)
}
