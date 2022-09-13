{{- with .get_events -}}
### `{{.command}}`

Retrieves the current, first and last event indices from a controller, and optionally 
fetches up to _count_ events starting from the current event index.

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
      "count": 3
    }
  }
}

{
  "message": {
    "reply": {
      "client-id": "QWERTY",
      "method": "get-events",
      "request-id": "AH173635G3",
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
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "7e6b0b90060593909bd148826e2e0aa03dfefaf1b2ab0316a15c87a36830b379"
}
```
{{end -}}


