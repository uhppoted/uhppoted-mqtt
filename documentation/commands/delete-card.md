### `delete-card`

Deletes a card from a controller.


```
Request:

topic: <root>/<requests>/device/card:delete

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "device-id": "uint32",
            "card-number": "uint32",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
device-id    (required) controller serial number
card-number  (required) card number
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "delete-card",
      "response": {
            "device-id": "<controller-id>",
            "card-number": "uint32",
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
card-number  card number
deleted      card delete success/fail
```


Example:
```
topic: uhppoted/gateway/requests/device/card:delete

{
  "message": {
    "request": {
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896,
      "card-number": 8165538
    }
  }
}

{
  "message": {
    "reply": {
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "method": "delete-card",
      "response": {
        "device-id": 405419896,
        "card-number": 8165538,
        "deleted": true
      }
    }
  }
}
```
