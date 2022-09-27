## Security

### HMAC

The first level of security is _HMAC_ and if configured as `required` in _uhppoted.conf_ requires that 
each request message be authenticated with an HMAC generated using a key associated with the request
`client-id` (the server response HMAC uses the server key). An HMAC simply provides an assurance that
the message has not been tampered with in transit (e.g. if it is passing through a broker not under your
control).

_uhppoted.conf_:
```
...
mqtt.security.HMAC.required = true
...
```

Request:
```
topic:  uhppoted/gateway/requests/devices:get

{
  "message": {
    "request": {
      "request-id": "AH173635G3",
      "client-id": "QWERTY",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896
    }
  },
  "hmac": "2574ee13c2a9aa1555a4200060e6250888a5c05c60897ee69b4a52347c102d9a"
}
```

Response:
```
{
  "message": {
    "reply": {
      "client-id": "QWERTY",
      "method": "get-device",
      "request-id": "AH173635G3",
      "response": {
        "device-id": 405419896,
        "device-type": "UTO311-L04",
        "ip-address": "192.168.1.100",
        "subnet-mask": "255.255.255.0",
        "gateway-address": "192.168.1.1",
        "mac-address": "00:12:23:34:45:56",
        "date": "2018-11-05",
        "version": "0892"
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "3fd19f56a23007ea702556938e1e91150fa211ebc4aca12f48df794362c9e9ce"
}

```

The response message will be authenticated with an HMAC generated using the server key.

### Authentication

The next level of security requires the request and response messages be digitally signed to provide
an assurance that they were actually originated by the client/server. _uhppoted-mqtt_ supports two mechanisms:

- a lightweight HOTP, where the client and server share the underlying HOTP key associated with the _client-id_ 
  of the request
- a more secure RSA digital signature where the server authenticates a request using the RSA public key associated
  with the _client-id_ of the request

Both of these mechanisms require an external means to securely exchange keys.

### HOTP

HOTP is a lightweight authentication mechanism based on a shared secret key and a counter that increases monontonically 
with each request. It is enabled in _uhppoted.conf_:
```
# MQTT
...
mqtt.security.authentication = HOTP
...
```

The keys are stored as `client-id::key` pairs in _/var/uhppoted/mqtt.hotp.secrets_, e.g.:
```
QWERTY      DAIOJ9BJQHPC7JBZ
```

and the corresponding counters are stored as `client-id::counter` pairs in _/var/uhppoted/mqtt.hotp.counters_, e.g.:
```
QWERTY      1093
```

Sample request/response:
```
{
  "message": {
    "request": {
      "request-id": "AH173635G3",
      "client-id": "QWERTY",
      "reply-to": "uhppoted/reply/97531",
      "nonce": 271,
      "hotp": "586787"
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
          ...
        }
      },
      "server-id": "uhppoted"
    }
  }
}
```

### RSA

RSA provides a stronger authentication mechanism based on digital signatures. The server uses the client public key
to verify each request and (optionally) signs each response with the server private key. RSA authentication is enabled
in _uhppoted.conf_:
```
# MQTT
...
mqtt.security.authentication = RSA
mqtt.security.outgoing.sign = true
...
```

The client public signing keys are stored as _\<client-id\>.pub_ PEM files in _/var/uhppoted/mqtt/rsa/signing_,
e.g.:
```
QWERTY.pub

-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAx8FbtmHSN8ui3eJN+CiM
dU1MEmHCzB9fGMplhnNjg/netI27ZVg+VvPMSvAF2c4Pq0MBYdhsOdU7i95SPRH4
...

```

The server signing key is stored as a PEM file in _/var/uhppoted/mqtt/rsa/signing/mqttd.key_, e.g.:
```
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDvmv0K0WQHN/wW
HgVxN5/adjhdMx0WKVWOFEVefN45/PjGIVOOKK80TS6Z/tJnIePD3tJRfi+gyI7D
...
```

Sample request/response:
```
{
  "message": {
    "signature": "VXLQgzQOHnjIFW6UFftWBYtdwluM3M7nbQD6fjLdSkuk/L8ahLfHsIEPCQF9ofkqEGaBG2Dl6QJtqYF825z8dLPsxbQA1bgMrdbpiVKiS09Vn4ubONIGmShQKcuoZuAzgsVeNbCsDW2MhSq/f6W/DUlKmD9PwgxMkzeKUCjM8bQ=",
    "request": {
      "request-id": "AH173635G3",
      "client-id": "QWERTY",
      "reply-to": "uhppoted/reply/97531",
      "nonce": 8
    }
  }
}
{
  "message": {
    "signature": "jVBJagU3KJKShK74RZBVbRGvv32/c3foq+6Dx98A4Pasic0iwNlpcG2fD4F3Zf3nZUsImhVdNTGRUtxzi7sbxI+fUR9STKCsRNDrjrP94gVotLk7GCT/mKyq58XSkHwluR6zj7P0qT3i9Y6U4Du5k8nhIzUObF9/Hff0WA6VtNVj9rhmIOg3pWJAhdjf+Hy6+9lxYjGwCc+3uZzo0wKaca68M3chONx78RlZvbXdr/S1AsvUx0avz+oX8lk4kGbJaq7g/FN9OSN4h8Hz6Al9/TYeUIMImfg3QgfPGujYd3tVSIfnYwmbYkEmQ4hJNQ8JkpDJ+zETf95Fo2Rt9yjJdw==",
    "reply": {
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "method": "get-devices",
      "nonce": 273,
      "request-id": "AH173635G3",
      "response": {
        "devices": {
          "201020304": {
            "device-type": "UTO311-L02",
            "ip-address": "192.168.1.101",
            "port": 60000
          },
          ...
        }
      }
    }
  },
  "hmac": "0e8e2b2f44e109fa9f80971b8aac0281b571d12efbb98aef14a16261bd53a09e"
}
```

### Authorisation

### Encryption

