### `record-special-events`

Enables/disables event logging for door open, close and pushbutton events.


```
Request:

topic: <root>/<requests>/device/special-events:set

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
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
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
topic: uhppoted/gateway/requests/device/special-events:set

{
  "message": {
    "request": {
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896,
      "enabled": true
    }
  }
}

{
  "message": {
    "reply": {
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "method": "record-special-events",
      "response": {
        "DeviceID": 405419896,
        "Enable": true,
        "Updated": true
      }
    }
  }
}
```
