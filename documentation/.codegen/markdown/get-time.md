{{- with .get_time -}}
### `{{.command}}`

Returns the controller date and time.

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
      "method": "get-time",
      "request-id": "AH173635G3",
      "response": {
        "date-time": "2022-09-08 11:01:04 PDT",
        "device-id": 405419896
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "b13b9bd7ddb5b84fde12bd19fb1cdd354084a287b43c414d109403e9bd6cc040"
}
```
{{end -}}


