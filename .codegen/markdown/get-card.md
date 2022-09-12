{{- with .get_cards -}}
### `{{.command}}`

Retrieves a list of the cards stored on a controller.

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
      "card-number": 8165538
    }
  }
}

{
  "message": {
    "reply": {
      "client-id": "QWERTY",
      "method": "get-card",
      "request-id": "AH173635G3",
      "response": {
        "device-id": 405419896,
        "card": {
          "card-number": 8165538,
          "start-date": "2021-01-01"
          "end-date": "2021-12-31",
          "doors": {
            "1": 1,
            "2": 0,
            "3": 0,
            "4": 1
          }
        }
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "61e3d8136ced7baca20440c2a5678319cce5dae57c8e8e7bed8cb484fa2a55d4"
}
```
{{end -}}
