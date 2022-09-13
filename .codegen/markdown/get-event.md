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
{{end -}}


