### `set-door-interlock`

Sets the door interlock mode for a controller.


```
Request:

topic: <root>/<requests>/device/door/interlock:set

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "device-id": "<controller-id>",
            "interlock": "<interlock>",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
device-id    (required) controller serial number
interlock    door interlock mode (0: none, 1:doors 1&2, 2:doors 3&4, 3:doors 1&2 and doors 3&4, 4:doors 1&2&3, 8:doors 1&2&3&4
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "set-door-interlock",
      "response": {
            "device-id": "<controller-id>",
            "interlock": "<interlock>",
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
device-id    controller serial number
interlock    door interlock mode (from request)
```


Example:
```
topic: uhppoted/gateway/requests/device/door/interlock:set

{
  "message": {
    "request": {
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896,
      "interlock": 4
    }
  }
}

{
  "message": {
    "reply": {
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "method": "set-door-interlokc",
      "response": {
        "device-id": 405419896,
        "interlock": 4
      }
    }
  }
}
```
