{{- with .get_device -}}
### `{{.command}}`

Returns the controller information for a UHPPOTE controller.

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
      "method": "get-device",
      "request-id": "AH173635G3",
      "response": {
        "address": {
          "IP": "192.168.1.100",
          "Port": 60000,
          "Zone": ""
        },
        "date": "2018-11-05",
        "device-id": 405419896,
        "device-type": "UTO311-L04",
        "gateway-address": "192.168.1.1",
        "ip-address": "192.168.1.100",
        "mac-address": "00:12:23:34:45:56",
        "subnet-mask": "255.255.255.0",
        "timezone": {},
        "version": "0892"
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "7e41dfbe8ed6eed02f8560a129d893b683909e04b40506d2b284cc22e8e4bb91"
}
```
{{end -}}


