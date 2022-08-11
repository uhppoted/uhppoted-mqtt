package mqtt

import (
	"fmt"
	"sync"
	"time"
)

const bucket_size = 20

type statistics struct {
	connected    []uint32
	disconnected []uint32
	subscribed   []uint32
	errors       []uint32

	N     int
	index int
	ch    <-chan time.Time

	sync.RWMutex
}

var stats = statistics{
	connected:    make([]uint32, bucket_size),
	disconnected: make([]uint32, bucket_size),
	subscribed:   make([]uint32, bucket_size),
	errors:       make([]uint32, bucket_size),

	N:     bucket_size,
	index: 0,
	ch:    time.Tick(15 * time.Second),
}

func init() {
	go func() {
		for {
			select {
			case <-stats.ch:
				(&stats).tick()
			}
		}
	}()
}

func (s *statistics) String() string {
	s.RLock()
	defer s.RUnlock()

	sum := func(b []uint32) uint64 {
		total := uint64(0)
		for _, v := range b {
			total += uint64(v)
		}

		return total
	}

	connected := sum(s.connected)
	disconnected := sum(s.disconnected)
	subscribed := sum(s.subscribed)
	errors := sum(s.errors)

	return fmt.Sprintf("connected:%v  disconnected:%v  subscribed:%v  errors:%v", connected, disconnected, subscribed, errors)
}

func (s *statistics) onConnected() {
	s.Lock()
	defer s.Unlock()

	s.connected[s.index]++
}

func (s *statistics) onDisconnected() {
	s.Lock()
	defer s.Unlock()

	s.disconnected[s.index]++
}

func (s *statistics) onSubscribed() {
	s.Lock()
	defer s.Unlock()

	s.subscribed[s.index]++
}

func (s *statistics) onError() {
	s.Lock()
	defer s.Unlock()

	s.errors[s.index]++
}

func (s *statistics) tick() {
	s.Lock()
	defer s.Unlock()

	index := (s.index + 1) % s.N

	s.connected[index] = 0
	s.disconnected[index] = 0
	s.subscribed[index] = 0
	s.errors[index] = 0
	s.index = index
}
