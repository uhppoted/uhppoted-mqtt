### `get-time-profiles`

Retrieves a range of time profiles from a controller.


```
Request:

topic: uhppoted/gateway/requests/device/time-profiles:get

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "device-id": "uint32",
            "from": "uint8",
            "to": "uint8",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
device-id    (required) controller serial number
from         (optional) start time profile ID [2..254]. Defaults to 2.
to           (optional) end time profile ID [2..254]. Defaults to 254.
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "get-time-profiles",
      "response": {
            "device-id": "uint32",
            "profiles": "array of record",
            "profile.id": "uint8",
            "profile.start-date": "date",
            "profile.end-date": "date",
            "profile.weekdays": "string list of weekday",
            "profile.segments": "array of time segments",
            "profile.segment.start": "time",
            "profile.segment.end": "time",
            "profile.linked-profile-id": "uint8",
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
device-id    controller serial number
profiles     array of time profile records
profile.id   time profile ID [2..254
profile.start-date time profile 'enabled from' date (inclusive)
profile.end-date time profile 'enabled until' date (inclusive)
profile.weekdays weekdays on which time profile is enabled
profile.segments time segments 1-3
profile.segment.start segment start time (HHmm)
profile.segment.end segment end time (HHmm)
profile.linked-profile-id (optional) ID of linked time profile [2..254]
```


Example:
```

  "message": {
    "request": {
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896,
      "from": 2,
      "to": 254
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
      "method": "get-time-profiles",
      "response": {
        "device-id": 405419896,
        "profiles": [
          {
            "end-date": "2021-12-31",
            "id": 29,
            "segments": [
              {
                "end": "17:00",
                "start": "08:30"
              },
              {
                "end": "00:00",
                "start": "00:00"
              },
              {
                "end": "00:00",
                "start": "00:00"
              }
            ],
            "start-date": "2021-01-01",
            "weekdays": "Monday"
          }
        ]
      }
    }
  }
}
```
