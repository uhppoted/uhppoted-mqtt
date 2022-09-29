# Command and Control Functions

This document describes the command and control functions implemented by _uhppoted_mqtt_.

## Security

Because it was originally designed as a gateway for mobile applications i.e. probably using a broker exposed
to the public and quite possibly using a public MQTT broker, an out-of-the-box installation of _uhppoted-mqtt_
has full authorisation, authentication and encryption enabled by default.

For the use case where it is running on a local network and using a trusted MQTT broker this is excessive
and unnecessary and also makes troubleshooting difficult. To disable all security, edit the
the following settings in the _uhppoted.conf_ file:

```
# MQTT
...
mqtt.security.HMAC.required = false
mqtt.security.authentication = NONE
mqtt.security.nonce.required = false
mqtt.security.outgoing.sign = false
mqtt.security.outgoing.encrypt = false
...
```

The security can then be increased incrementally as required. The details are described in:

- [HMAC](https://github.com/uhppoted/uhppoted-mqtt/blob/master/documentation/commands/security.md#hmac)
- [authentication](https://github.com/uhppoted/uhppoted-mqtt/blob/master/documentation/commands/security.md#authentication)
- [authentication](https://github.com/uhppoted/uhppoted-mqtt/blob/master/documentation/commands/security.md#authentication)
    - [HOTP](https://github.com/uhppoted/uhppoted-mqtt/blob/master/documentation/commands/security.md#hotp)
    - [RSA](https://github.com/uhppoted/uhppoted-mqtt/blob/master/documentation/commands/security.md#rsa)
- [authorisation](https://github.com/uhppoted/uhppoted-mqtt/blob/master/documentation/commands/security.md#authorisation)
- [encryption](https://github.com/uhppoted/uhppoted-mqtt/blob/master/documentation/commands/security.md#encryption)
- [nonce](https://github.com/uhppoted/uhppoted-mqtt/blob/master/documentation/commands/security.md#nonce)

## Command structure

Each command comprises:
- topic
- request message
- response message

e.g. for the _get-device_ command:
```
topic: uhppoted/gateway/requests/device:get

request:
   {
     "message": {
       "request": {
         "client-id": "QWERTY",
         "request-id": "AH173635G3",
         "reply-to": "uhppoted/reply/97531",
         "device-id": 405419896
       }
     }
   }

response:
   {
     "message": {
       "reply": {
         "server-id": "uhppoted"
         "client-id": "QWERTY",
         "request-id": "AH173635G3",
         "method": "get-device",
         "response": {
           "device-id": 405419896,
           "device-type": "UTO311-L04",
           "ip-address": "192.168.1.100",
           "subnet-mask": "255.255.255.0",
           "gateway-address": "192.168.1.1",
           "mac-address": "00:12:23:34:45:56",
           "date": "2018-11-05",
           "version": "0892"
         }
       }
     }
   }

```

### Topic

The message _topic_ can be thought of as the MQTT equivalent of a REST URL and is used internally to dispatch the 
request to the correct handler. Each command has a unique topic which is prefixed by the 'root' topic
and the 'requests' section defined in _uhppoted.conf_:
```
...
# MQTT
mqtt.topic.root = uhppoted/gateway
mqtt.topic.requests = ./requests
...
```

e.g. for the _get-device_ command:
```
uhppoted/gateway/requests/device:get
```

### Request

A _request_ message comprises the following fields:

- `client-id`
- `request-id`
- `reply-to`

followed by the request parameters that are unique to each command.

The `client-id` is required and identifies the requesting application for purposes of authentication and authorisation,
using the permissions granted in:
- \<etc/uhppoted\>/mqtt.permissions.users 
- \<etc/uhppoted\>/mqtt.permissions.groups

The `request-id` is an optional client identifier for the request that is echoed back in the response.

The `reply-to` is an optional topic for the MQTT response message to simplify internal dispatch within a client
application. If not provided, it defaults to the `<root>/<replies>/<command>` topics defined in _uhppoted.conf_:
```
...
# MQTT
mqtt.connection.client.key = /usr/local/etc/com.github.uhppoted/mqtt/client.key
mqtt.topic.root = uhppoted/gateway
mqtt.topic.requests = ./requests
mqtt.topic.replies = ./replies
mqtt.topic.events = ./events
mqtt.topic.system = ./system
...
```

e.g. for the _get-device_ command:
```
   {
     "message": {
       "request": {
         "client-id": "QWERTY",
         "request-id": "AH173635G3",
         "reply-to": "uhppoted/reply/97531",
         "device-id": 405419896
       }
     }
   }
```

### Response

A _response_ message comprises the following fields:

- `client-id`
- `server-id`
- `request-id`
- `method`

followed by the response values that are unique to each command.

The `server-id` identifies the responding server - mostly for the use case where a client might be interacting with
multiple servers but also as a system/security check that the response originated from the expected source.

The `client-id` is echoed from the `client-id` from the originating request for use with dispatching 
received replies to the appropriated handler.

The `request-id` is echoed from the `request-id` of the originating request to allow the client to match the 
asynchronous reply with the originating request.

The `method` field identifies the originating command for the response. It can also be used to dispatch the response
to the appropriate handler for clients that do not use the `request-id` field for dispatch.

e.g. for the _get-device_ command:
```
{
  "message": {
    "reply": {
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "method": "get-device",
      "response": {
        "device-id": 405419896,
        "device-type": "UTO311-L04",
        "ip-address": "192.168.1.100",
        "subnet-mask": "255.255.255.0",
        "gateway-address": "192.168.1.1",
        "mac-address": "00:12:23:34:45:56",
        "date": "2018-11-05",
        "version": "0892"
      }
    }
  }
}

```

## List of Commands

- [`get-devices`](get-devices.md)

- [`get-device`](get-device.md)

- [`get-time`](get-time.md)

- [`set-time`](set-time.md)

- [`get-door-delay`](get-door-delay.md)

- [`set-door-delay`](set-door-delay.md)

- [`get-door-control`](get-door-control.md)

- [`set-door-control`](set-door-control.md)

- [`record-special-events`](record-special-events.md)

- [`open-door`](open-door.md)

- [`get-status`](get-status.md)

- [`get-cards`](get-cards.md)

- [`get-card`](get-card.md)

- [`put-card`](put-card.md)

- [`get-events`](get-events.md)

- [`get-event`](get-event.md)

- [`delete-card`](delete-card.md)

- [`delete-cards`](delete-cards.md)

- [`get-time-profile`](get-time-profile.md)

- [`set-time-profile`](set-time-profile.md)

- [`clear-time-profiles`](clear-time-profiles.md)

- [`get-time-profiles`](get-time-profiles.md)

- [`set-time-profiles`](set-time-profiles.md)

- [`set-task-list`](set-task-list.md)

- [`acl-show`](acl-show.md)

- [`acl-grant`](acl-grant.md)

- [`acl-revoke`](acl-revoke.md)

- [`acl-upload-file`](acl-upload-file.md)

- [`acl-upload-s3`](acl-upload-s3.md)

- [`acl-upload-http`](acl-upload-http.md)

- [`acl-download-file`](acl-download-file.md)

- [`acl-download-s3`](acl-download-s3.md)

- [`acl-download-http`](acl-download-http.md)

- [`acl-compare-file`](acl-compare-file.md)

- [`acl-compare-s3`](acl-compare-s3.md)

- [`acl-compare-http`](acl-compare-http.md)

