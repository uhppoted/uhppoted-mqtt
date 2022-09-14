{{- with .get_event -}}
### `{{.command}}`

Retrieves an event from the controller.

{{template "request"  . -}}
{{template "response" . }}

Example:
```
{
  "message": {
    "request": {
{{- template "request-preamble"}}
      "device-id": 405419896,
      "event-index": 50
    }
  }
}

{
  "message": {
    "reply": {
{{- template "response-preamble"}}
      "method": "get-event",
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
      }
    }
  }
}
```
{{end -}}


