## v0.7.x

## IN PROGRESS

- [ ] `set-time-profiles`
- [ ] `get-acl` with time-profiles
- [ ] `set-acl` with time-profiles
- [ ] `show` with time-profiles
- [ ] `grant` with time-profiles

- [x] `clear-time-profiles`
- [x] `get-time-profiles`
- [x] `set-time-profile`
- [x] `get-time-profile`
- [x] Update to latest `uhppote-core` and `uhppoted-api`
- [x] Add handling for `nil` events in response to `get-status`

## TODO

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