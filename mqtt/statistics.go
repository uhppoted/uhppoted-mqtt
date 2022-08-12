package mqtt

import (
	"time"

	"github.com/uhppoted/uhppoted-mqtt/log"
)

type statistics struct {
	disconnects []uint32
	max         uint64

	index        int
	disconnected chan uint32
	tick         <-chan time.Time
}

var stats = statistics{
	disconnects: make([]uint32, 20),
	max:         10,

	index:        0,
	disconnected: make(chan uint32),
	tick:         time.Tick(15 * time.Second),
}

func init() {
	stats.monitor()
}

func (s *statistics) onDisconnected() {
	s.disconnected <- uint32(1)
}

func (s *statistics) monitor() {
	sum := func(b []uint32) uint64 {
		total := uint64(0)
		for _, v := range b {
			total += uint64(v)
		}

		return total
	}

	go func() {
		for {
			select {
			case N := <-stats.disconnected:
				stats.disconnects[stats.index] += N
				count := sum(stats.disconnects)
				log.Infof(LOG_TAG, "DISCONNECT COUNT:%v of %v", count, s.max)
				if count >= stats.max {
					log.Fatalf(LOG_TAG, "DISCONNECT COUNT %v REACHED MAXIMUM ALLOWED (%v)", count, stats.max)
				}

			case <-stats.tick:
				index := (stats.index + 1) % len(stats.disconnects)
				stats.disconnects[index] = 0
				stats.index = index
			}
		}
	}()
}
