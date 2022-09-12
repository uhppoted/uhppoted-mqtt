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
      "device-id": 405419896
    }
  }
}

{
  "message": {
    "reply": {
      "client-id": "QWERTY",
      "method": "get-cards",
      "request-id": "AH173635G3",
      "response": {
        "device-id": 405419896,
        "cards": [
          8165537,
          8165539,
          8165538
        ]
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "fae3fcfc4f63f960b7e8c0e3bda283d77757ed95025de5af9f65baa5851ef814"
}
```
{{end -}}


