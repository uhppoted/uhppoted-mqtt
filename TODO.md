# TODO

### IN PROGRESS

- [ ] Greengrass + Paho:
```
x509: certificate signed by unknown authority 
```
      - https://github.com/google/go-github/issues/1049
      - https://groups.google.com/g/golang-nuts/c/v5ShM8R7Tdc
      - https://stackoverflow.com/questions/62828165/got-x-509-certificate-signed-by-unknown-authority-when-the-server-tried-to-sen
      - https://medium.com/the-go-journey/x509-certificate-signed-by-unknown-authority-running-a-go-app-inside-a-docker-container-a12869337eb

- [ ] Figure out missing conf.LockfileRemove for Linux
- [x] Remove ARM7 specific files after removing softlock

- [ ] Write up AWS GreenGrass setup
      - [x] Create separate policy/group/user for Greengrass CLI
      - [x] `core` and `thing` provisioning
      - [x] uhppoted-mqtt installation and configuration
      - [ ] Discover API/script
      - https://docs.aws.amazon.com/greengrass/v2/developerguide/greengrass-discover-api.html
      - https://iot.stackexchange.com/questions/6347/connecting-cellular-module-sim7070g-to-aws-mqtt/6350
      - [NodeRed/GreenGrass](https://iot.stackexchange.com/questions/2646/deploy-scripts-to-aws-greengrass-without-aws-lambda)


- [x] Remove _false_ `from mqtt.acl.verify`

## TODO

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

## NOTES

1. [os_arch.go](https://gist.github.com/camabeh/a02e6846e00251e1820c784516c0318f)
