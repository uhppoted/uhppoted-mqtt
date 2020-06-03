## v0.6.x

## IN PROGRESS

- [x] acl::show
- [x] acl::grant
- [x] inject response metadata in MQTTD
- [x] acl::revoke
- [x] acl::upload
- [x] acl::download
- [x] acl::compare
- [x] Add report URL to response from acl:compare
- [x] Move devices to ACL struct (or UHPPOTED ?)
- [x] Remove context from handlers
- [x] Make report from acl:Compare a bit less download-y
- [ ] move device commands to separate package
- [ ] rework devices commands to match ACL
- [ ] send `Error` reply as msgError
- [ ] rethink embedding reply in body
- [ ] unit test for GetCardByIndex with return code 0xffffffff

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
