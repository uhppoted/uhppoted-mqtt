### `set-antipassback`

Sets the controller antipassback mode:

| Anti-passback | Description                                                  |
|---------------|--------------------------------------------------------------|
| _disabled_    | No anti-passback                                             |
| _(1:2);(3:4)_ | Doors 1 and 2 are interlocked, doors 3 and 4 are interlocked |
| _(1,3):(2,4)_ | Doors 1 and 3 are interlocked with doors 2 and 4             |
| _1:(2,3)_     | Door 1 is interlocked with doors 2 and 3                     |
| _1:(2,3,4)_   | Door 1 is interlocked with doors 2,3 and 4                   |

where _interlocked_ means a card will be swiped through a second time on a door until it has 
been swiped through at the _interlocked_ door. e.g: if the anti-passback mode is _(1,3):(2,4),
a card swiped through at either of doors 1 or 3 will be denied access at doors 1 and 3 until 
it has been swiped through at either of doors 2 or 4. Likewise a card swiped through at either
of doors 2 or 4 will be denied access at doors 2 and 4 until is has been swiped through at 
either of doors 1 or 3.


```
Request:

topic: <root>/<requests>/device/antipassback:set

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "device-id": "<controller-id>",
            "anti-passback": "<anti-passback>",
        }
    }
}

request-id     (optional) message ID, returned in the response
client-id      (required) client ID for authentication and authorisation (if enabled)
reply-to       (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                          configured reply topic) if not provided.
device-id      (required) controller serial number
anti-passback  (required) anti-passback mode. One of:
               - disabled
               - (1:2);(3:4)
               - (1,3):(2,4)
               - 1:(2,3)
               - 1:(2,3,4)
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "set-time",
      "response": {
            "device-id": "<controller-id>",
            "anti-passback": "<anti-passback>",
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
device-id    controller serial number
anti-passback   (required) anti-passback mode. One of:
              - disabled
              - (1:2);(3:4)
              - (1,3):(2,4)
              - 1:(2,3)
              - 1:(2,3,4)
```


Example:
```
topic: uhppoted/gateway/requests/device/time:set

{
  "message": {
    "request": {
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896,
      "anti-passback": "(1,3):(2,4)"
    }
  }
}

{
  "message": {
    "reply": {
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "method": "set-antipassback",
      "response": {
        "device-id": 405419896,
        "anti-passback": "(1,3):(2,4)"
      }
    }
  }
}
```
