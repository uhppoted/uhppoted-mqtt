package mqtt

import (
	"time"

	"github.com/uhppoted/uhppoted-mqtt/log"
)

type statistics struct {
	disconnects []uint32
	interval    time.Duration
	max         uint32

	disconnected chan uint32
	tick         <-chan time.Time
}

var stats = statistics{
	disconnects: make([]uint32, 60),
	interval:    5 * time.Minute,
	max:         10,

	disconnected: make(chan uint32),
	tick:         time.Tick(1 * time.Second),
}

func init() {
	stats.monitor()
}

func SetDisconnectsInterval(interval time.Duration) {
	if interval > 60*time.Second {
		stats.interval = interval
	} else {
		stats.interval = 60 * time.Second
	}
}

func SetMaxDisconnects(N uint32) {
	stats.max = N
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
		start := time.Now()
		index := 0

		for {
			select {
			case N := <-stats.disconnected:
				stats.disconnects[index] += N
				count := sum(stats.disconnects)
				log.Infof(LOG_TAG, "DISCONNECT %v of %v in %v", count, s.max, s.interval)
				if count >= uint64(stats.max) {
					log.Fatalf(LOG_TAG, "DISCONNECT COUNT %v REACHED MAXIMUM ALLOWED (%v)", count, stats.max)
				}

			case <-stats.tick:
				N := len(stats.disconnects)
				step := stats.interval / time.Duration(N)
				delta := time.Now().Sub(start)
				bucket := float64(delta) / float64(step)
				next := int(bucket) % 60

				for next != index {
					index = (index + 1) % 60
					stats.disconnects[index] = 0
				}
			}
		}
	}()
}
