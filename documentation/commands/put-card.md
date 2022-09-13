### `put-card`

Adds or updates a card record on a controller.


```
Request:

topic: uhppoted/gateway/requests/device/card:put

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "device-id": "uint32",
            "card": "record",
            "card-number": "uint32",
            "start-date": "date",
            "end-date": "date",
            "doors": "{1:uint8, 2:uint8, 3:uint8, 4:uint8}",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
device-id    (required) controller serial number
card         card record
card-number  card number
start-date   card 'valid from' date (inclusive)
end-date     card 'valid until' date (inclusive)
doors        door [1..4] access rights
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "put-card",
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
{
  "message": {
    "request": {
      "request-id": "AH173635G3",
      "client-id": "QWERTY",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896,
      "card": {
        "card-number": 8165538,
        "start-date": "2021-01-01",
        "end-date": "2021-12-31",
        "doors": {
          "1": true,
          "2": false,
          "3": 55,
          "4": false
        }
      }
    }
  }
}

{
  "message": {
    "reply": {
      "client-id": "QWERTY",
      "method": "put-card",
      "request-id": "AH173635G3",
      "response": {
        "device-id": 405419896
        "card": {
          "card-number": 8165538,
          "start-date": "2021-01-01"
          "end-date": "2021-12-31",
          "doors": {
            "1": true,
            "2": false,
            "3": 55,
            "4": false
          },
        },
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "535122d6cac2b5ce29f294fdb9739bc90a39e8ab7913ccb8bd010571ecf51013"
}
```
