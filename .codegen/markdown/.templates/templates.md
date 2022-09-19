{{define "request"}}
```
Request:

topic: <root>/<requests>/{{ .request.topic }}

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
{{- range .request.fields}}
            "{{.field}}": "{{.value}}",
{{- end}}
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
{{- range .request.fields}}
{{printf "%-12s" .field}} {{.description}}
{{- end}}
```
{{end}}


{{define "request-preamble"}}
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
{{- end}}


{{define "response"}}
```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "{{.command}}",
      "response": {
{{- range .response.fields}}
            "{{.field}}": "{{.value}}",
{{- end}}
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
{{- range .response.fields}}
{{printf "%-12s" .field}} {{.description}}
{{- end}}
```
{{end}}

{{define "response-preamble"}}
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
{{- end}}



