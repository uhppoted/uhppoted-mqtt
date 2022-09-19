### `acl-revoke`

Revokes system-wide access control permissions for a card. The permissions are
removed from any existing permissions.


```
Request:

topic: <root>/<requests>/acl/card:revoke

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
      "method": "acl-revoke",
      "response": {
            "revoked": "bool",
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
revoked      revoke success/fail
```


Example:
```
topic: uhppoted/gateway/requests/acl/card:revoke

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
        "Dungeon"
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
        "revoked": true
      }
    }
  }
}
```
