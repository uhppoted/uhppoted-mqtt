{{- with .get_devices -}}
### `{{.command}}`

Returns a list of all UHPPOTE controllers found via a UDP broadcast on the local LAN or specifically
configured in _uhppoted.conf_.

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
      "nonce": 5
    }
  }
}

{
  "message": {
    "reply": {
      "client-id": "QWERTY",
      "method": "get-devices",
      "request-id": "AH173635G3",
      "response": {
        "devices": {
          "201020304": {
            "device-type": "UTO311-L02",
            "ip-address": "192.168.1.101",
            "port": 60000
          },
          "303986753": {
            "device-type": "UTO311-L03",
            "ip-address": "192.168.1.100",
            "port": 60000
          },
          "405419896": {
            "device-type": "UTO311-L04",
            "ip-address": "192.168.1.100",
            "port": 60000
          }
        }
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "77b8917b80e18837fc9c13e451194d17f2712fb95b8f5283f9c762dfe1ed4f55"
}
```
{{end -}}


