### `acl-upload`

Extracts an access control list from a set of controllers and stores it to the URL in the request. The URL
can be any URL that will accept an file, including:

- `file://<file path>` (e.g. _file:///var/uhppoted/uploaded/ACL.tar.gz_)
- `http://<path>` (e.g. _http://localhost:8080/uploaded/ACL.tar.gz_)
- `https://<path>` (e.g. _https://localhost:8080/uploaded/ACL.tar.gz_)
- `s3://<bucket>/<object>` (e.g. _s3://uhppoted/uploaded/ACL.tar.gz_)

Uploads to Amazon S3 expect the credentials and bucket information to be configured in _uhppoted.conf_:
```
aws.credentials = /etc/uhppoted/aws.credentials
; aws.profile = default
; aws.region = us-east-1
```

_aws.credentials_:
```
[default]
aws_access_key_id = AKIAQMNWIVKBYA57IRWH
aws_access_key_id = AKIZ.............QYV
aws_secret_access_key = FRE................................zuyqt

```


```
Request:

topic: <root>/<requests>/acl/acl:upload

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "url": "URL",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
url          Destination URL for the ACL file
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "acl-upload",
      "response": {
            "url": "URL",
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
url          URL of uploaded file
```


Example:
```
topic: uhppoted/gateway/requests/acl/acl:upload

{
  "message": {
    "request": {
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
    "request": {
      "url": "file:///var/uhppoted/uploaded/ACL.tar.gz",
    }
  }
}

{
  "message": {
    "reply": {
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "method": "acl:upload",
      "response": {
        "uploaded": "file:///var/uhppoted/uploaded/uhppoted.tar.gz"
      }
    }
  }
}
```
