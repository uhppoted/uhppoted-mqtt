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

Authorisation restricts the operations which a client has permission. It is appropriate for systems that are accessed
by a number of users/clients, not all of whom should have full access to the system. 

Authorisation is based on:
- the `client-id` from the request message
- user group membership
- group permissions

1. For authorisation to be effective, the client-id should be authenticated i.e. an authentication method (HOTP and/or RSA)
should be enabled in _uhppoted.conf_:
```
mqtt.security.authentication = HOTP, RSA
```

2. User group membership is defined in the _`<etc`>/mqtt.permissions.users_ file and specify the groups to which 
   a user belongs, e.g. for a _mqtt.permissions.users_ file that looks like:
```
QWERTY   system, admin, user
UIOP     user
```

- QWERTY is a member of the _system_, _admin_ and _user_ groups
- UIOP is a member of only the _user_ group

3. Permissions for each group are defined in the _`<etc`>/mqtt.permissions.groups_ file and specify the topics
   a member of the group is allowed to access, e.g.:
```
system    *:*
admin     events:*, event:*, acl:card:*
user      acl/card:show
```

grants the following permissions:

- _system_ has unrestricted access to all topics
- _admin_ can only see events and add/delete/update cards
- _user_ can only see card permissions


### Encryption

Encrytion is intended to protect message content in transit - it is horrifically complicated, seldom necessary and
only really intended for systems that use public MQTT brokers.

Encryption is configured in _uhppoted.conf_:
```
...
mqtt.security.rsa.keys = /etc/uhppoted/mqtt/rsa
mqtt.security.outgoing.encrypt = true
...
```

- `mqtt.security.rsa.keys` sets the folder for RSA keys which contains both the encryption and signing keys. 
   Encryption keys are then located in the _.../encryption_ subfolder (e.g.  `/etc/uhppoted/mqtt/rsa/encryption`).
   Under the encyption keys subfolder:

   - client public keys are stored in PEM format as _<client-id>.pub_, e.g.:
```
   /etc/uhppoted/mqtt/rsa/encryption/QWERTY.pub

-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDFt51kdvOYNUgs8noX5TutPxgg
....
....
70yIjpbXzXEhLRZpPQIDAQAB
-----END PUBLIC KEY-----
```

   - the private key used to encrypt outgoing responses is stored in _.../mqttd.key_
   - the private key used to encrypt outgoing events is stored in _.../events.key_


Requests can be encrypted on a message-by-message basis - an encrypted request will include:
- a `key`
- a `signature`
- an encrypted `request`

e.g:
```
get-device-encrypted:
  mqtt publish --topic 'uhppoted/gateway/requests/device:get'
               --message '{ 
                             "message": { 
                                "key": "LRtq7KaKvsCP8VvaRXRsoc2+R5T...WmI3nSQ==",
                                "signature": "G/3cEtzhZ+5iy.....zyahv8=",
                                "request": "EJo5lNjfYl....NJUCQ==" 
                            },
                             "hmac": "5a0d5ffd....8bba27" 
                          }'
```

- the `key` is a base-64 string corresponding to a generated AES-256 key encrypted by the client
  using its RSA private key. It should be unique for every request (i.e. do not reuse keys)
- the `signature` is the RSA signature of the request
- the request is a base-64 string corresponding to the request encrypted under the AES `key` using CBC chaining.

An encrypted response is similar:
```
{
  "message": {
    "key": "leFH4a2zYIawJrcBtxNNK1D46dlyaAtqV+JuFdI8OeoQrihUVh9/Ul2Qslti7Ow49bQolTVFS6ww3DWYby3qjZHx2sfwGr1yOnWSOXrUVo/bff237iscl+OuSVT24jlDKKtW0Njiq48+3S2DCDJ9h2nW/vl5Q76TtGZXqIQAHpw=",
    "reply": "SQecW39OR0ILcY4ErW/WdBlwo+J7TNvuGG+KORTRm1ljxkR9LXe60Lx2YJXlQJIRxW8yl//vHIQ+s6tbUedue2lKREPMT6iF5emzHr3BmT4OipdLkp+ogD+RZSDIWKhNSl7V29G6sGxYB5xf4akjDH8ylLVi/KXD/0e07pwygsE0uLifjRWqibZb9m3og7VFO1NCnaSXNeiVZcdK2T5sP/vT7UfgouJvm2goSvTQ6yiF2Mh9xVpTKL799d/jxUl9MHgCSNzIP6lZAdS1rbq9p39UhjGwsFM6WndqKrndrDvKeSAVYAXm+Q6LTCfUrM+b4MYcgPdn4FinE+bnt+m9Enwqchk3OyE1j0e9+s4rnBt53FwNJF+3LqWmKlXCwjju7oLG8bbH6tCGxxGN/Yktgwnz1mLMeDNLkBMoQsTiAZjxzh1QSKTSrhk6WS7dkwTmGJKK14OaNyMsTWx3GHDCzoCQjqWwvLxK54VGJiIMNeNH8Asn8uzxvwRWt6hpSD6U0U3CPNSyTLWwpZPVPc71v/4up/ySvcJLKGIUcs6D3pTa4EgVjO3MNOpB2Bdu/knm"
  },
  "hmac": "0b020617d250503647b4bf6e3d395c7f938a673b64967fb9f274fcae677e9693"
}
```

- the `key` is a base-64 string corresponding to a generated AES-256 key encrypted by the server using the `mqttd.key` 
  RSA private key. It is unique for every request.
- the `signature` is the RSA signature of the reply, using the server signing key
- the reply is a base-64 string corresponding to the reply encrypted under the AES `key` using CBC chaining.

### Nonce

The above security enhancements still leave the protocol vulnerable to replay attacks i.e. where a request such as `open-door`
can be recorded and resent later by an unauthorised client. To mitigate this scenario, `uhppoted-mqtt` optionally mandates
a `nonce` for every request and response - the `nonce` being a monotonically increasing number that cannot be reused. 

A client request with a `nonce` that is not greater than the `nonce` of the last received request from that client is silently 
discarded as _invalid_. The `nonce` in a reply is the server `nonce` and is likewise increases monotonically - it has less
security value than the client `nonce` but can be used as some protection against fake responses.

The `nonce` functionality is enabled or disabled in _uhppoted.conf_:
```
...
mqtt.security.nonce.required = true
mqtt.security.nonce.server = /var/com.github.uhppoted/mqtt.nonce
mqtt.security.nonce.clients = /var/com.github.uhppoted/mqtt.nonce.counters
...
```

and requires that client authentication is both enabled and required. If enabled, each request must contain a unique
(and increasing) `nonce` in the message:
```
{
  "message": {
    "request": {
      "request-id": "AH173635G3",
      "client-id": "QWERTY",
      "reply-to": "uhppoted/reply/97531",
      "nonce": 271,
      ...
    }
  }
}
```

Notes:
- the server `nonce` value is stored in the _<var>/mqtt.nonce_ file
```
mqttd                 376
```

- the client `nonce` values are stored in the _<var>/mqtt.nonce.counters_ file
```
...
QWERTY                173
...

```

