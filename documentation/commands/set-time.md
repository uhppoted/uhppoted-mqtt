### `get-time`

Sets the controller date and time.


```
Request:

topic: uhppoted/gateway/requests/device/time:get

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "device-id": "<controller-id>",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the configured reply topic) if not provided.
device-id    (required) controller serial number
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
{
  "message": {
    "request": {
      "request-id": "AH173635G3",
      "client-id": "QWERTY",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896,
      "date-time": "2022-09-09 11:25:04"
    }
  }
}

{
  "message": {
    "reply": {
      "client-id": "QWERTY",
      "method": "set-time",
      "request-id": "AH173635G3",
      "response": {
        "device-id": 405419896,
        "date-time": "2022-09-09 11:25:04 PDT"
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "1e3a06eb2022315bc5fb6335d4ff3ef62c75c4a6c5b395fd43e5f1876703bb27"
}
```
