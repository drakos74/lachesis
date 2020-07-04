package network

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// WorldClock defines a unit of time or processing cycle
type WorldClock struct {
	tick      chan struct{}
	tock      chan Event
	eventPool *Events
	cycles    int
}

func (wc WorldClock) startTicking() {
	for range wc.tick {
		wc.cycles++
		// leave some time to warm up, and use the same amount to move to the next events
		if wc.cycles > wc.eventPool.warmUp {
			idx := wc.eventPool.index
			if idx < len(wc.eventPool.events) {
				log.Info().
					Str("Type", "EVENT").
					Msg(fmt.Sprintf("apply new event at %d - %d = %v", wc.cycles, idx, wc.eventPool.events[idx]))
				event := wc.eventPool.events[idx]
				wc.tock <- event
				wc.eventPool.index++
			}
			wc.cycles = 0
		}
	}
}
