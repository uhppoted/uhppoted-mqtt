### `set-door-passcodes`

Sets up to four supervisor passcodes for a controller door, with valid passcodes being in the range [1..999999]. 
Invalid passcodes are replaces by 0 (no code).


```
Request:

topic: <root>/<requests>/device/door/passcodes:set

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "device-id": "<controller-id>",
            "door": "<door-id>",
            "passcodes": "<passcodes>",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
device-id    (required) controller serial number
door         (required) door (1..4)
passcodes    array of passcodes ([1..999999]). Only the first four entries are used.
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "set-door-passcodes",
      "response": {
            "device-id": "<controller-id>",
            "door": "<door-id>",
            "passcodes": "<passcodes>",
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
device-id    controller serial number
door         door (1..4) from the request
passcodes    passcodes assigned to door
```


Example:
```
topic: uhppoted/gateway/requests/device/door/passcodes:set

{
  "message": {
    "request": {
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896,
      "door": 3,
      "passcodes": [12345,999999,54321]
    }
  }
}

{
  "message": {
    "reply": {
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "method": "set-door-control",
      "response": {
        "device-id": 405419896,
        "door": 3,
        "passcodes": [12345,999999,54321]
      }
    }
  }
}
```
