### `set-time-profiles`

Stores a list of time profiles to a controller.


```
Request:

topic: <root>/<requests>/device/time-profiles:set

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
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
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
device-id    (required) controller serial number
profiles     array of time profile record
profile.id   time profile ID [2..254
profile.start-date time profile 'enabled from' date (inclusive)
profile.end-date time profile 'enabled until' date (inclusive)
profile.weekdays weekdays on which time profile is enabled
profile.segments time segments 1-3
profile.segment.start segment start time (HHmm)
profile.segment.end segment end time (HHmm)
profile.linked-profile-id (optional) ID of linked time profile [2..254]
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "set-time-profiles",
      "response": {
            "device-id": "uint32",
            "warnings": "array of string",
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
device-id    controller serial number
warnings     list of warning messages
```


Example:
```
topic: uhppoted/gateway/requests/device/time-profiles:set

  "message": {
    "request": {
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896,
      "profiles": [
        {
          "id": 29,
          "start-date": "2021-01-01",
          "end-date": "2021-12-31",
          "weekdays": "Monday,Wednesday,Thursday",
          "segments": [
            {
              "start": "08:15",
              "end": "11:30"
            },
            {
              "start": "14:05",
              "end": "17:45"
            }
          ],
          "linked-profile": 30
        },
        {
          "id": 31,
          "start-date": "2021-01-01",
          "end-date": "2021-12-31",
          "weekdays": "Monday,Wednesday,Thursday",
          "segments": [
            {
              "start": "08:15",
              "end": "11:30"
            },
            {
              "start": "14:05",
              "end": "17:45"
            }
          ],
          "linked-profile": 32
        },
        {
          "id": 32,
          "start-date": "2021-01-01",
          "end-date": "2021-12-31",
          "weekdays": "Monday,Wednesday",
          "segments": [
            {
              "start": "08:30",
              "end": "15:30"
            }
          ],
          "linked-profile": 33
        },
        {
          "id": 33,
          "start-date": "2021-01-01",
          "end-date": "2021-12-31",
          "weekdays": "Saturday,Sunday",
          "segments": [
            {
              "start": "10:30",
              "end": "17:30"
            }
          ]
        }
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
      "method": "get-time-profile",
      "response": {
      "method": "get-time-profiles",
      "response": {
        "device-id": 405419896,
        "warnings": [
          "profile 29 : linked time profile 30 is not defined"
        ]
      }
    }
  }
}
```
