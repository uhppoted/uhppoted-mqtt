# Command and Control Functions

*STATUS: IN PROGRESS*

This document describes the command and control functions implemented by _uhppoted_mqtt_.

## Security

Because it was originally designed as a gateway for mobile applications possibly using a public
MQTT broker, an out-of-the-box installation of _uhppoted-mqtt_ has full authorisation, authentication
and encryption enabled.

For the use case where it is running on a local network and using a trusted MQTT broker this is excessive
and unnecessary and also makes troubleshooting difficult. To disable all security, edit the
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

{{ range .commands}}
- [`{{.command}}`]({{.command}}.md)
{{end}}
