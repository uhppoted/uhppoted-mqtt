# TODO

### IN PROGRESS

- [ ] Write up AWS GreenGrass setup
      - [x] Fix all document links so that they are relative
      - (?) Create _UhppotedGreengrassTokenExchangeRoleAccess_ policy
```
Encountered error - User: arn:aws:iam::NNNN:user/uhppoted-greengrass is not authorized to perform: iam:GetPolicy on resource: policy arn:aws:iam::aws:policy/UhppotedGreengrassTokenExchangeRoleAccess because no identity-based policy allows the iam:GetPolicy action (Service: Iam, Status Code: 403, Request ID: 428d25c4-7bbb-42b9-bf8d-71c001d0f60e); No permissions to lookup managed policy, looking for a user defined policy...
IAM policy named "UhppotedGreengrassTokenExchangeRoleAccess" already exists. Please attach it to the IAM role if not already
Configuring Nucleus with provisioned resource details...
```

      - [ ] Discover API/script
      - https://docs.aws.amazon.com/greengrass/v2/developerguide/greengrass-discover-api.html
      - https://iot.stackexchange.com/questions/6347/connecting-cellular-module-sim7070g-to-aws-mqtt/6350
      - [NodeRed/GreenGrass](https://iot.stackexchange.com/questions/2646/deploy-scripts-to-aws-greengrass-without-aws-lambda)

- [ ] Remove startup warnings for missing encryption/signing/etc files if auth is not enabled.
- [ ] Clean up Paho logging
- [ ] MQTT v5

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
