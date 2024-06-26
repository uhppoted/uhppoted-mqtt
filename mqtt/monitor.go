package mqtt

import (
	"sync"
	"time"

	"github.com/uhppoted/uhppoted-lib/monitoring"
	"github.com/uhppoted/uhppoted-mqtt/log"
)

type SystemMonitor struct {
	mqttd *MQTTD
}

var alive = sync.Map{}

func NewSystemMonitor(mqttd *MQTTD) *SystemMonitor {
	return &SystemMonitor{
		mqttd: mqttd,
	}
}

func (m *SystemMonitor) Alive(monitor monitoring.Monitor, msg string) error {
	event := struct {
		Alive struct {
			SubSystem string `json:"subsystem"`
			Message   string `json:"message"`
		} `json:"alive"`
	}{
		Alive: struct {
			SubSystem string `json:"subsystem"`
			Message   string `json:"message"`
		}{
			SubSystem: monitor.ID(),
			Message:   msg,
		},
	}

	now := time.Now()
	last, ok := alive.Load(monitor.ID())
	interval := 60 * time.Second

	if ok && time.Since(last.(time.Time)).Round(time.Second) < interval {
		return nil
	}

	if err := m.mqttd.send(&m.mqttd.Encryption.SystemKeyID, m.mqttd.Topics.System, nil, event, msgSystem, false); err != nil {
		log.Warnf("monitoring", "%v", err)
		return err
	}

	alive.Store(monitor.ID(), now)

	return nil
}

func (m *SystemMonitor) Alert(monitor monitoring.Monitor, msg string) error {
	event := struct {
		Alert struct {
			SubSystem string `json:"subsystem"`
			Message   string `json:"message"`
		} `json:"alert"`
	}{
		Alert: struct {
			SubSystem string `json:"subsystem"`
			Message   string `json:"message"`
		}{
			SubSystem: monitor.ID(),
			Message:   msg,
		},
	}

	if err := m.mqttd.send(&m.mqttd.Encryption.SystemKeyID, m.mqttd.Topics.System, nil, event, msgSystem, true); err != nil {
		log.Warnf("monitoring", "%v", err)
		return err
	}

	return nil
}
