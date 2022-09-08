# Command and Control Functions

*STATUS: IN PROGRESS*

This document describes the command and control functions implemented by _uhppoted_mqtt_.

## Security

Beacuse it was originally designed as a gateway for mobile applications possibly using a public
MQTT broker, an out-of-the-box installation of _uhppoted-mqtt_ has full authorisation, authentication
and encryption enabled.

For the use case where it is running on a local network and using a trusted MQTT broker this is excessive
and unnecessary and also makes initial troubleshooting difficult. To disable all security, edit the
the following settings in the _uhppoted.conf_ file:

```
# MQTT
...
mqtt.security.authentication = NONE
mqtt.security.HMAC.required = false
mqtt.security.nonce.required = false
mqtt.security.outgoing.sign = false
mqtt.security.outgoing.encrypt = false
...
```

The security can then be increased incrementally as required.

### Authentication

### Authorisation

### Encryption

## Commands

- `[get-devices]`(get-devices.md)- `[get-device]`(get-device.md)- `[get-time]`(get-time.md)- `set-time`
- `get-door-delay`
- `set-door-delay`
- `get-door-control`
- `set-door-control`
- `record-special-events`
- `get-status`
- `get-cards`
- `delete-cards`
- `get-card`
- `put-card`
- `delete-card`
- `get-time-profile`
- `set-time-profile`
- `clear-time-profiles`
- `get-time-profiles`
- `set-time-profiles`
- `set-task-list`
- `get-events`
- `get-event`
- `[open-door]`(open-door.md)- `acl-show`
- `acl-grant`
- `acl-revoke`
- `acl-upload-file`
- `acl-upload-s3`
- `acl-upload-http`
- `acl-download-file`
- `acl-download-s3`
- `acl-download-http`
- `acl-compare-file`
- `acl-compare-s3`
- `acl-compare-http`

