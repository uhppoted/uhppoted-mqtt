### `record-special-events`

Enables/disables event logging for door open, close and pushbutton events.


```
Request:

topic: uhppoted/gateway/requests/device/special-events:set

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "device-id": "<controller-id>",
            "enabled": "<bool>",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the configured reply topic) if not provided.
device-id    (required) controller serial number
enabled      (required) true/false
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "record-special-events",
      "response": {
            "device-id": "<controller-id>",
            "door": "<door-id>",
            "control": "<mode>",
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
device-id    controller serial number
door         door (1..4) from the request
control      door control mode (normally open, normally closed or controlled)
```


Example:
```
{
  "message": {
    "request": {
      "request-id": "AH173635G3",
      "client-id": "QWERTY",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896,
      "enabled": true
    }
  }
}

{
  "message": {
    "reply": {
      "client-id": "QWERTY",
      "method": "record-special-events",
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
