### `set-time`

Sets the controller date and time.


```
Request:

topic: <root>/<requests>/device/time:set

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "device-id": "<controller-id>",
            "date-time": "YYYY-mm-dd hh:mm:ss"
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
device-id    (required) controller serial number
date-time    (required) set controller system date and time in this format: YYYY-mm-dd hh:mm:ss
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "get-time",
      "response": {
            "device-id": "<controller-id>",
            "date-time": "<datetime>",
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
device-id    controller serial number
date-time    controller system date and time
```


Example:
```
topic: uhppoted/gateway/requests/device/time:get

{
  "message": {
    "request": {
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896,
      "date-time": "2022-09-09 11:25:04"
    }
  }
}

{
  "message": {
    "reply": {
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "method": "set-time",
      "response": {
        "device-id": 405419896,
        "date-time": "2022-09-09 11:25:04 PDT"
      }
    }
  }
}
```
