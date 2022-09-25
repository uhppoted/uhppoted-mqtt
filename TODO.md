# TODO

### IN PROGRESS

- [x] Added ACL.NoVerify to configuration
- [ ] Add support for downloading non-tar.gz files

- [ ] Protocol documentation
      - [ ] Preamble
      - [ ] SUMMARY.md

- [ ] `set-time-profiles` is returning incorrect method in response
- [ ] Ewww - fix record-special-events response
```
{
  "message": {
    "reply": {
      "client-id": "QWERTY54",
      "method": "record-special-events",
      "nonce": 230,
      "request-id": "AH173635G3",
      "response": {
        "DeviceID": 405419896,
        "Enable": true,
        "Updated": true
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "c7767d44fd50c7cda2ee5e1bf1af2d86b8a2367a4a330fda0a84427cc8c1c8ef"
}
```

- [ ] Write up AWS GreenGrass setup
      - [ ] Create separate policy/group/user for Greengrass CLI
      - https://github.com/awsdocs/aws-iot-greengrass-v2-developer-guide/issues/20
      - https://docs.aws.amazon.com/greengrass/v2/developerguide/greengrass-discover-api.html

- (?) Generate uhppoted.conf if it doesn't exist

## TODO

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

## NOTES

1. [os_arch.go](https://gist.github.com/camabeh/a02e6846e00251e1820c784516c0318f)
