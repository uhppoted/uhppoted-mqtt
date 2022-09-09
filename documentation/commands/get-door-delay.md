### `get-door-delay`

Retrieves the open delay for a controller door.


```
Request:

topic: uhppoted/gateway/requests/device/door/delay:get

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "device-id": "<controller-id>",
            "door": "<door-id>",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the configured reply topic) if not provided.
device-id    (required) controller serial number
door         (required) door (1..4)
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "get-door-delay",
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
delay        door open delay
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
      "door": 3
    }
  }
}

{
  "message": {
    "reply": {
      "client-id": "QWERTY",
      "method": "get-door-delay",
      "request-id": "AH173635G3",
      "response": {
        "device-id": 405419896,
        "door": 3,
        "delay": 7
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "506fe61c581fbd3cd438aa89d984b6f382bdf95a8f0170a1b3a2edd073fe4c03"
}
```
