package ticker

import (
	"testing"
	"time"
)

func TestRestart(t *testing.T) {
	ticker := New(100*time.Millisecond)
	time.Sleep(50*time.Millisecond)

	select {
	case <-ticker.C:
		t.Errorf("Got a tick too early(#0)")
	default:
	}

	time.Sleep(60*time.Millisecond)
	select {
	case <-ticker.C:
	default:
		t.Errorf("Did not get a tick, but expected to (#0)")
	}

	ticker.Stop() // this's here to test how Restart() will work if the ticker is already stopped

	ticker.Restart(200 * time.Millisecond)
	time.Sleep(110*time.Millisecond)
	select {
	case <-ticker.C:
		t.Errorf("Got a tick too early (#1)")
	default:
	}

	time.Sleep(100*time.Millisecond)
	select {
	case <-ticker.C:
	default:
		t.Errorf("Did not get a tick, but expected to (#1)")
	}
}
