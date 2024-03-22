# TODO

### IN PROGRESS

- [ ] [Prebuilt Docker containers](https://github.com/uhppoted/uhppoted/issues/47)
       - [x] docker-dev
       - [x] docker-ghcr
       - [x] compose.yml
       - [x] README
       - [x] CHANGELOG
       - [ ] Update uhppote to use uhppoted-rest docker
       - [ ] `Incorrect listener address`

- [x] [`restore-default-parameters`](https://github.com/uhppoted/uhppoted/issues/48)

- [ ] Remove startup warnings for missing encryption/signing/etc files if auth is not enabled.
- [ ] Clean up Paho logging
- [ ] MQTT v5

## TODO

- [ ] [Sparkplug B](https://github.com/eclipse-sparkplug/sparkplug)
- [ ] [MQTT Dash](https://iot.stackexchange.com/questions/6561/generic-mobile-applications-for-smart-home-devices)
- [ ] [UCANs](https://ucan.xyz/)
- [ ] (optionally) Generate uhppoted.conf if it doesn't exist
- [ ] Make reconnect time configurable
- [ ] Relook at encoding reply content - maybe json.RawMessage can preserve the field order
- [ ] Replace values passed in Context with initialised struct
- [ ] last-will-and-testament (?)
- [ ] publish add/delete card, etc to event stream
- [ ] MQTT v5.0
- [ ] [JSON-RPC](https://en.wikipedia.org/wiki/JSON-RPC) (?)
- [ ] Add to CLI
- [ ] Non-ephemeral key transport:  https://tools.ietf.org/html/rfc5990#appendix-A
- [ ] user:open/get permissions require matching card number 
- [ ] [AEAD](http://alexander.holbreich.org/message-authentication)
- [ ] Support for multiple brokers
- [ ] NACL/tweetnacl
- [ ] Report system events for e.g. listen bound/not bound

### Documentation

- [ ] TeX protocol description
- [ ] godoc
- [ ] build documentation
- [ ] install documentation
- [ ] user manuals
- [ ] man/info page

### Other

1.  github project page
2.  Integration tests
3.  Verify fields in listen events/status replies against SDK:
    - battery status can be (at least) 0x00, 0x01 and 0x04
4.  EventLogger 
    - MacOS: use [system logging](https://developer.apple.com/documentation/os/logging)
    - Windows: event logging
5.  Update file watchers to fsnotify when that is merged into the standard library (1.4 ?)
    - https://github.com/golang/go/issues/4068
6. [Teserakt E2E encryption](https://teserakt.io)
7. [Fernet encryption](https://asecuritysite.com/encryption/fernet)
8. [IoT standards](https://iot.stackexchange.com/questions/5363/mqtt-json-format-for-process-automation-industry)
9. [StackExchange: MQTT security tests](https://iot.stackexchange.com/questions/452/what-simple-security-tests-can-i-perform-on-my-mqtt-network)
10. [VerneMQ](https://vernemq.com)
11. [SparkplugB](https://cogentdatahub.com/connect/mqtt/sparkplug-b)

## NOTES

1. [os_arch.go](https://gist.github.com/camabeh/a02e6846e00251e1820c784516c0318f)
