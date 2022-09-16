### `acl-grant`

Grants system-wide access control permissions for a card. The permissions are
added to any existing permissions.


```
Request:

topic: uhppoted/gateway/requests/acl/card:grant

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "card-number": "uint32",
            "start-date": "date",
            "end-date": "date",
            "doors": "array of string",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
card-number  card number
start-date   card 'valid from' date (inclusive)
end-date     card 'valid until' date (inclusive)
doors        string list of door names
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "acl-grant",
      "response": {
            "granted": "bool",
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
granted      grant success/fail
```


Example:
```
{
  "message": {
    "request": {
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
    "request": {
      "card-number": 8165538,
      "start-date": "2022-01-01",
      "end-date": "2022-12-31",
      "doors": [
        "Gryffindor",
        "Slytherin"
      ]
    }
  }
}

{
  "message": {
    "reply": {
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "method": "acl:grant",
      "response": {
        "granted": true
      }
    }
  }
}
```
