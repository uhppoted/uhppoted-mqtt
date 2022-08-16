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
	disconnects: make([]uint32, 20),
	interval:    5 * time.Minute,
	max:         10,

	disconnected: make(chan uint32),
	tick:         time.Tick(TICK),
}

const TICK = 1 * time.Second

func init() {
	stats.monitor()
}

func SetDisconnectsInterval(interval time.Duration) {
	stats.interval = interval
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
		last := 0

		for {
			select {
			case N := <-stats.disconnected:
				stats.disconnects[0] += N
				count := sum(stats.disconnects)
				log.Infof(LOG_TAG, "DISCONNECT %v of %v in %v", count, s.max, s.interval)
				if count >= uint64(stats.max) {
					log.Fatalf(LOG_TAG, "DISCONNECT COUNT %v REACHED MAXIMUM ALLOWED (%v)", count, stats.max)
				}

			case <-stats.tick:
				N := len(stats.disconnects)
				step := stats.interval / time.Duration(N)
				delta := time.Now().Sub(start)
				bucket := int(float64(delta) / float64(step))

				// TODO rework this as ring buffer
				for bucket > last {
					disconnects := append([]uint32{0}, stats.disconnects...)
					stats.disconnects = disconnects[0:N]
					last += 1
				}
			}
		}
	}()
}
