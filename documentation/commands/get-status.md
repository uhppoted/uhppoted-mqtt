### `get-status`

Retrieves the controller status.


```
Request:

topic: uhppoted/gateway/requests/device/status:get

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "device-id": "uint32",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
device-id    (required) controller serial number
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "get-status",
      "response": {
            "device-id": "uint32",
            "status": "record",
            "system-datetime": "datetime",
            "door-states": "{ 1:bool, 2:bool, 3:bool, 4:bool }",
            "door-buttons": "{ 1:bool, 2:bool, 3:bool, 4:bool }",
            "input-state": "uint8",
            "relay-state": "uint8",
            "system-error": "uint8",
            "special-info": "uint8",
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
status       controller status record
system-datetime controller system date/time
door-states  door open/closed states
door-buttons door button states
input-state  input state bitset
relay-state  relay state bitset
system-error system error code
special-info internal state information code
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
      "device-id": 405419896
    }
  }
}

{
  "message": {
    "reply": {
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "method": "get-status",
      "response": {
        "device-id": 405419896,
        "status": {
          "door-buttons": {
            "1": false,
            "2": false,
            "3": false,
            "4": false
          },
          "door-states": {
            "1": false,
            "2": false,
            "3": false,
            "4": false
          },
          "event": {
            "access-granted": true,
            "card-number": 8165537,
            "device-id": 0,
            "direction": 1,
            "direction-text": "in",
            "door-id": 1,
            "event-id": 69,
            "event-reason": 44,
            "event-reason-text": "remote open door",
            "event-type": 2,
            "event-type-text": "door",
            "timestamp": "2021-08-10 10:28:32 PDT"
          },
          "input-state": 0,
          "relay-state": 0,
          "sequence-id": 0,
          "special-info": 0,
          "system-datetime": "2022-09-12 10:47:31 PDT",
          "system-error": 0
        }
      }
    }
  }
}
```
