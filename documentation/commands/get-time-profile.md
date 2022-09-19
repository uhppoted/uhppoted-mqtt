### `get-time-profile`

Retrieves a time profile from a controller.


```
Request:

topic: <root>/<requests>/device/time-profile:get

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "device-id": "uint32",
            "profile-id": "uint8",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
device-id    (required) controller serial number
profile-id   (required) time profile ID [2..254]
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "get-time-profile",
      "response": {
            "device-id": "uint32",
            "time-profile": "record",
            "id": "uint8",
            "start-date": "date",
            "end-date": "date",
            "weekdays": "string list of weekday",
            "segments": "array of time segments",
            "segment.start": "time",
            "segment.end": "time",
            "linked-profile-id": "uint8",
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
device-id    controller serial number
time-profile time profile record
id           time profile ID [2..254
start-date   time profile 'enabled from' date (inclusive)
end-date     time profile 'enabled until' date (inclusive)
weekdays     weekdays on which time profile is enabled
segments     time segments 1-3
segment.start segment start time (HHmm)
segment.end  segment end time (HHmm)
linked-profile-id (optional) ID of linked time profile [2..254]
```


Example:
```
topic: uhppoted/gateway/requests/device/time-profile:get

  "message": {
    "request": {
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896,
      "profile-id": 29
    }
  }
}

{
  "message": {
    "reply": {
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "method": "get-time-profile",
      "response": {
        "device-id": 405419896,
        "time-profile": {
          "id": 29,
          "start-date": "2021-01-01",
          "end-date": "2021-12-31",
          "weekdays": "Monday"
          "segments": [
            {
              "start": "08:30",
              "end": "17:00"
            },
            {
              "start": "00:00",
              "end": "00:00"
            },
            {
              "start": "00:00",
              "end": "00:00"
            }
          ],
        }
      }
    }
  }
}
```
