{{- with .get_status -}}
### `{{.command}}`

Retrieves the controller status.

{{template "request"  . -}}
{{template "response" . }}

Example:
```
topic: uhppoted/gateway/requests/{{ .request.topic }}

{
  "message": {
    "request": {
{{- template "request-preamble"}}
      "device-id": 405419896
    }
  }
}

{
  "message": {
    "reply": {
{{- template "response-preamble"}}
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
{{end -}}


