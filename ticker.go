package ticker

// This ticker supports changing the interval while running

import (
	"time"
	"sync/atomic"
)

const (
	status_stopped = iota
	status_running
	status_runned
	status_stopping
)

type ticker struct {
	ticker *time.Ticker
	C chan time.Time
	stopSignal chan struct{}
	lastTriggerTime time.Time
	status uint32
}

func New(interval time.Duration) *ticker {
	result := &ticker{}
	result.C = make(chan time.Time)
	result.stopSignal = make(chan struct{})
	result.status = status_stopped
	result.start(interval)
	return result
}

func (ticker *ticker) start(interval time.Duration) {
	if !atomic.CompareAndSwapUint32(&ticker.status, status_stopped, status_running) {
		return
	}
	ticker.lastTriggerTime = time.Now()
	ticker.ticker = time.NewTicker(interval)
	ticker.run()
	atomic.StoreUint32(&ticker.status, status_runned)
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
	if !atomic.CompareAndSwapUint32(&ticker.status, status_runned, status_stopping) {
		return
	}
	ticker.flushChan()
	ticker.stopSignal <- struct{}{}
	ticker.ticker.Stop()
	atomic.StoreUint32(&ticker.status, status_stopped)
}

func (ticker *ticker) Stop() {
	ticker.stop()
}

func (ticker *ticker) Restart(newInterval time.Duration) {
	ticker.stop()
	ticker.start(newInterval)
}
