### `get-time`

Returns the controller date and time.


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
      "device-id": 405419896
    }
  }
}

{
  "message": {
    "reply": {
      "client-id": "QWERTY",
      "method": "get-time",
      "request-id": "AH173635G3",
      "response": {
        "date-time": "2022-09-08 11:01:04 PDT",
        "device-id": 405419896
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "b13b9bd7ddb5b84fde12bd19fb1cdd354084a287b43c414d109403e9bd6cc040"
}
```
