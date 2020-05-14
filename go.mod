module github.com/uhppoted/uhppoted-mqtt

go 1.14

require (
	github.com/eclipse/paho.mqtt.golang v1.2.1-0.20200121105743-0d940dd29fd2
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/uhppoted/uhppote-core v0.6.1
	github.com/uhppoted/uhppoted-api v0.6.1
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b // indirect
	golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae
)

replace(
	github.com/uhppoted/uhppote-core => ../uhppote-core
)
