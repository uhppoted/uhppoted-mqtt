### `get-cards`

Retrieves a list of the cards stored on a controller.


```
Request:

topic: uhppoted/gateway/requests/device/cards:get

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
      "method": "get-cards",
      "response": {
            "device-id": "<controller-id>",
            "cards": "[]uint32",
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
device-id    controller serial number
cards        list of card numbers
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
      "method": "get-cards",
      "request-id": "AH173635G3",
      "response": {
        "device-id": 405419896,
        "cards": [
          8165537,
          8165539,
          8165538
        ]
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "fae3fcfc4f63f960b7e8c0e3bda283d77757ed95025de5af9f65baa5851ef814"
}
```
