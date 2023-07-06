### `activate-keypads`

Activates/deactivates the reader access keypads for a controller.


```
Request:

topic: <root>/<requests>/device/door/keypads:set

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "device-id": "<controller-id>",
            "keypads": "map[uint8]bool",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
device-id    (required) controller serial number
keypads      map of activated readers (unlisted readers are deactivated)
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "activate-keypads",
      "response": {
            "device-id": "<controller-id>",
            "keypads": "map[uint8]bool",
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
device-id    controller serial number
keypads      map of readers to activated status
```


Example:
```
topic: uhppoted/gateway/requests/device/door/keypads:set

{
  "message": {
    "request": {
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896,
      "keypads": {
        "1": true,
        "2": true,
        "3": false,
        "4": true
      }
    }
  }
}

{
  "message": {
    "reply": {
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "method": "set-door-keypads",
      "response": {
        "device-id": 405419896,
        "keypads": {
          "1": true,
          "2": true,
          "3": false,
          "4": true
        }
      }
    }
  }
}
```
