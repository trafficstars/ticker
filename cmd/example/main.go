package main

import (
	"log"
	"time"

	"github.com/trafficstars/ticker"
)

func main() {
	ticker := ticker.New(10*time.Millisecond)
	time.Sleep(5*time.Millisecond)

	select {
	case <-ticker.C:
		log.Printf("Got a tick too early(#0)")
	default:
	}

	time.Sleep(6*time.Millisecond)
	select {
	case <-ticker.C:
	default:
		log.Printf("Did not get a tick, but expected to (#0)")
	}

	ticker.Restart(20 * time.Millisecond)
	time.Sleep(11*time.Millisecond)
	select {
	case <-ticker.C:
		log.Printf("Got a tick too early (#1)")
	default:
	}

	time.Sleep(10*time.Millisecond)
	select {
	case <-ticker.C:
	default:
		log.Printf("Did not get a tick, but expected to (#1)")
	}
}
