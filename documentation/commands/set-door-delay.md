### `set-door-delay`

Sets the open delay for a door on a controller.


```
Request:

topic: <root>/<requests>/device/door/delay:set

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
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
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
topic: uhppoted/gateway/requests/device/door/delay:set

{
  "message": {
    "request": {
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896,
      "door": 3,
      "delay": 8
    }
  }
}

{
  "message": {
    "reply": {
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "method": "set-door-delay",
      "response": {
        "device-id": 405419896,
        "door": 3,
        "delay": 8
      }
    }
  }
}
```
