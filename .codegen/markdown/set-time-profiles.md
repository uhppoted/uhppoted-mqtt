{{- with .set_time_profiles -}}
### `{{.command}}`

Stores a list of time profiles to a controller.

{{template "request"  . -}}
{{template "response" . }}

Example:
```
  "message": {
    "request": {
{{- template "request-preamble"}}
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
{{- template "response-preamble"}}
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
{{end -}}
