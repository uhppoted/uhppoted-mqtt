### `clear-time-profiles`

Clears all time profiles on a a controller.


```
Request:

topic: <root>/<requests>/device/time-profiles:delete

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "device-id": "uint32",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
device-id    (required) controller serial number
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "clear-time-profiles",
      "response": {
            "device-id": "<controller-id>",
            "deleted": "bool",
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
device-id    controller serial number
deleted      clear time profiles success/fail
```


Example:
```
topic: uhppoted/gateway/requests/device/time-profiles:delete

  "message": {
    "request": {
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896
    }
  }
}

{
  "message": {
    "reply": {
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "method": "clear-time-profiles",
      "response": {
        "cleared": true,
        "device-id": 405419896
      },
    }
  },
}
```
