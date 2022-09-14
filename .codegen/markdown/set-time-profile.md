{{- with .set_time_profile -}}
### `{{.command}}`

Adds or updates a time profile on a controller.

{{template "request"  . -}}
{{template "response" . }}

Example:
```
{
  "message": {
    "request": {
{{- template "request-preamble"}}
      "device-id": 405419896,
      "profile": {
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
        "linked-profile": 3
      }
    }
  }
}

{
  "message": {
    "reply": {
{{- template "response-preamble"}}
      "method": "set-time-profile",
      "response": {
        "device-id": 405419896,
        "time-profile": {
          "id": 29,
          "start-date": "2021-01-01",
          "end-date": "2021-12-31",
          "weekdays": "Monday,Wednesday,Thursday"
          "segments": [
            {
              "end": "11:30",
              "start": "08:15"
            },
            {
              "end": "17:45",
              "start": "14:05"
            },
            {
              "end": "00:00",
              "start": "00:00"
            }
          ],
          "linked-profile": 3
        }
      }
    }
  }
}
```
{{end -}}
