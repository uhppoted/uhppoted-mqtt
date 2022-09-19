### `get-card`

Retrieves a card record from a controller.


```
Request:

topic: <root>/<requests>/device/card:get

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
      "method": "get-card",
      "response": {
            "device-id": "<controller-id>",
            "card": "record",
            "card-number": "uint32",
            "start-date": "date",
            "end-date": "date",
            "doors": "{1:uint8, 2:uint8, 3:uint8, 4:uint8}",
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
device-id    controller serial number
card         card record
card-number  card number
start-date   card 'valid from' date (inclusive)
end-date     card 'valid until' date (inclusive)
doors        door [1..4] access rights
```


Example:
```
topic: uhppoted/gateway/requests/device/card:get

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
      "method": "get-card",
      "response": {
        "device-id": 405419896,
        "card": {
          "card-number": 8165538,
          "start-date": "2021-01-01"
          "end-date": "2021-12-31",
          "doors": {
            "1": 1,
            "2": 0,
            "3": 0,
            "4": 1
          }
        }
      }
    }
  }
}
```
