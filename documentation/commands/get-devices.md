### `get-devices`

Returns a list of all UHPPOTE controllers found via a UDP broadcast on the local LAN or specifically
configured in _uhppoted.conf_.


```
Request:

topic: uhppoted/gateway/requests/devices:get

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "get-devices",
      "response": {
            "device-id": "<controller-id>",
            "device-type": "<controller-type>",
            "ip-address": "<address>",
            "port": "<UDPv4 port>",
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
device-id    controller serial number
device-type  controller type (UTO311-L0x)
ip-address   controller IPv4 address
port         UDP port for controller commands
```


Example:
```
{
  "message": {
    "request": {
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
    }
  }
}

{
  "message": {
    "reply": {
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "method": "get-devices",
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
      }
    }
  }
}
```
