### `get-events`

Retrieves the current, first and last event indices from a controller, and optionally 
fetches up to _count_ events starting from the current event index.


```
Request:

topic: uhppoted/gateway/requests/device/events:get

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "device-id": "uint32",
            "count": "uint32",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
device-id    (required) controller serial number
count        (optional) number of events to fetch (defaults to 0)
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "get-events",
      "response": {
            "device-id": "<controller-id>",
            "current": "uint32",
            "first": "uint32",
            "last": "uint32",
            "events": "array of event",
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
current      current event index
first        first event index
last         last event index
events       list of up to 'count' events
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
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896,
      "count": 3
    }
  }
}

{
  "message": {
    "reply": {
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "method": "get-events",
      "response": {
        "device-id": 405419896,
        "current": 26,
        "first": 1,
        "last": 69,
        "events": [
          {
            "access-granted": true,
            "card-number": 8165533,
            "device-id": 405419896,
            "direction": 1,
            "direction-text": "in",
            "door-id": 1,
            "event-id": 24,
            "event-reason": 0,
            "event-reason-text": "",
            "event-type": 2,
            "event-type-text": "door",
            "timestamp": "2019-07-24 20:31:24 PDT"
          },
          {
            "access-granted": true,
            "card-number": 8165534,
            "device-id": 405419896,
            "direction": 1,
            "direction-text": "in",
            "door-id": 4,
            "event-id": 25,
            "event-reason": 0,
            "event-reason-text": "",
            "event-type": 2,
            "event-type-text": "door",
            "timestamp": "2019-07-31 20:04:00 PDT"
          },
          {
            "access-granted": false,
            "card-number": 8165535,
            "device-id": 405419896,
            "direction": 0,
            "direction-text": "",
            "door-id": 4,
            "event-id": 26,
            "event-reason": 0,
            "event-reason-text": "",
            "event-type": 1,
            "event-type-text": "le swipe",
            "timestamp": "2019-07-31 20:04:32 PDT"
          }
        ]
      }
    }
  }
}
```
