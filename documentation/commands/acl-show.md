### `acl-show`

Retrieves the system-wide access control permissions for a card.


```
Request:

topic: uhppoted/gateway/requests/acl/card:show

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "card-number": "uint32",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
card-number  (required) card number
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "acl-show",
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
{
  "message": {
    "request": {
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
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
      "method": "acl:show",
      "response": {
        "card-number": 8165538,
        "permissions": [
          {
            "door": "Gryffindor",
            "end-date": "2021-12-31",
            "start-date": "2021-01-01"
          },
          {
            "door": "Slytherin",
            "end-date": "2021-12-31",
            "start-date": "2021-01-01"
          }
        ]
      }
    }
  }
}
```
