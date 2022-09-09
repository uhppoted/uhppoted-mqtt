### `set-door-delay`

Sets the control mode for a controller door.


```
Request:

topic: uhppoted/gateway/requests/device/door/delay:set

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "device-id": "<controller-id>",
            "door": "<door-id>",
            "delay": "<delay>",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the configured reply topic) if not provided.
device-id    (required) controller serial number
door         (required) door (1..4) to open
delay        door open delay
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "set-door-delay",
      "response": {
            "device-id": "<controller-id>",
            "door": "<door-id>",
            "delay": "<delay>",
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
delay        door open duration
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
      "door": 3,
      "control": "normally closed"
    }
  }
}

{
  "message": {
    "reply": {
      "client-id": "QWERTY",
      "method": "set-door-control",
      "request-id": "AH173635G3",
      "response": {
        "device-id": 405419896,
        "door": 3,
        "control": "normally closed"
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "5d2550d6d5a8c11c79fc531a4391fe80f6a66f241792b8432ee102ff0c20c833"
}
```
