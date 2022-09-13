### `get-event`

Retrieves an event from the controller.


```
Request:

topic: uhppoted/gateway/requests/device/event:get

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "device-id": "uint32",
            "event-index": "uint32",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
device-id    (required) controller serial number
event-index  (required) index of event to fetch
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "get-event",
      "response": {
            "device-id": "uint32",
            "event": "event record",
            "event.timestamp": "datetime",
            "event.device-id": "uint32",
            "event-id": "uint32",
            "event.type": "uint8",
            "event.type-text": "string",
            "event.door-id": "uint8",
            "event.card-number": "uint32",
            "event.access-granted": "bool",
            "event.direction": "uint8",
            "event.direction-text": "string",
            "event.reason": "uint8",
            "event.reason-text": "string",
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
device-id    controller serial number
event        event details
event.timestamp last event date/time
event.device-id last event controller serial number
event-id     last event index
event.type   last event type code
event.type-text last event type
event.door-id last event door
event.card-number last event card number
event.access-granted last event access granted
event.direction last event direction (1:IN, 2:OUT)
event.direction-text last event direction ('in', 'out')
event.reason last event reason code
event.reason-text last event reason
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
      "event-index": 50
    }
  }
}

{
  "message": {
    "reply": {
      "client-id": "QWERTY54",
      "method": "get-event",
      "request-id": "AH173635G3",
      "response": {
        "device-id": 405419896,
        "event": {
          "access-granted": true,
          "card-number": 8165538,
          "device-id": 405419896,
          "direction": 1,
          "direction-text": "in",
          "door-id": 4,
          "event-id": 50,
          "event-reason": 0,
          "event-reason-text": "",
          "event-type": 2,
          "event-type-text": "door",
          "timestamp": "2019-08-09 16:18:55 PDT"
        }
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "2079c940984153774d4338a0ffb3e99d66a368e0237b99249463c4a47236767a"
}
```
