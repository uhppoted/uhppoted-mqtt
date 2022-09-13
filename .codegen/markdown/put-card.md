{{- with .put_card -}}
### `{{.command}}`

Adds or updates a card record on a controller.

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
      "card": {
        "card-number": 8165538,
        "start-date": "2021-01-01",
        "end-date": "2021-12-31",
        "doors": {
          "1": true,
          "2": false,
          "3": 55,
          "4": false
        }
      }
    }
  }
}

{
  "message": {
    "reply": {
      "client-id": "QWERTY",
      "method": "put-card",
      "request-id": "AH173635G3",
      "response": {
        "device-id": 405419896
        "card": {
          "card-number": 8165538,
          "start-date": "2021-01-01"
          "end-date": "2021-12-31",
          "doors": {
            "1": true,
            "2": false,
            "3": 55,
            "4": false
          },
        },
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "535122d6cac2b5ce29f294fdb9739bc90a39e8ab7913ccb8bd010571ecf51013"
}
```
{{end -}}
