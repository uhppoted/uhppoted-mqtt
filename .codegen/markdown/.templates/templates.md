{{define "request"}}
```
Request:

topic: {{ .request.topic }}

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
