# FAQ

1. My MQTT broker is on an secure network - how do I disable signing and encryption and HMAC and ..?

In _uhppoted.conf_:
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

2. Why am I still getting old events even though I have set `mqtt.alerts.retained = false`

The MQTT broker probably still has retained events from when `mqtt.alerts.retained` was set to `false`. You
need to manually clear the events - normally by sending an empty retained message to the topics:

```
mqtt publish --topic 'uhppoted/gateway/events' --message '' -r
mqtt publish --topic 'uhppoted/gateway/system' --message '' -r
```