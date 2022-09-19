### `acl-compare`

Compares the access control permissions on a set of controllers with a file retrieved from a URL. The URL
can be any URL that will return a file, including:

- `file://<file path>` (e.g. _file:///var/uhppoted/ACL.tar.gz_)
- `http://<path>` (e.g. _http://localhost:8080/ACL.tar.gz_)
- `https://<path>` (e.g. _https://localhost:8080/ACL.tar.gz_)
- `s3://<bucket>/<file>` (e.g. _s3://uhppoted/ACL.tar.gz_)

Downloads from Amazon S3 expect the credentials and bucket information to be configured in _uhppoted.conf_:
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

topic: <root>/<requests>/acl/acl:compare

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "acl": "URL",
            "report": "URL",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
acl          Source URL for the ACL file
report       Destination URL for the report file
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "acl-compare",
      "response": {
            "url": "URL",
            "report": "list of record",
            "report.controller": "uint32",
            "report.controller.diffent": "uint32",
            "report.controller.extraneous": "uint32",
            "report.controller.updated": "uint32",
            "report.controller.unchanged": "uint32",
            "report.controller.missing": "uint32",
            "report.controller.unchanged": "uint32",
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
url          URL for the uploaded report file
report       list of the changes made to each controller
report.controller controller ID
report.controller.diffent number of cards on the controller that have the same card number but different permissions
report.controller.extraneous number of cards on the controller that are not in the ACL file
report.controller.updated number of cards updated on controller
report.controller.unchanged number of cards unchanged on controller
report.controller.missing number of card in the ACL file that are not present on the controller
report.controller.unchanged number of cards on the controller that match the ACL file
```


Example:
```
topic: uhppoted/gateway/requests/acl/acl:compare

{
  "message": {
    "request": {
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
    "request": {
      "acl": "file:///var/uhppoted/ACL.tar.gz",
      "report": "file:///var/uhppoted/report.tar.gz"
    }
  }
}

{
  "message": {
    "reply": {
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "method": "acl:compare",
      "response": {
        "url": "file:///var/uhppoted/report.tar.gz"
        "report": {
          "201020304": {
            "different": 0,
            "extraneous": 0,
            "missing": 0,
            "unchanged": 3
          },
          "303986753": {
            "different": 0,
            "extraneous": 0,
            "missing": 2,
            "unchanged": 1
          },
          "405419896": {
            "different": 0,
            "extraneous": 0,
            "missing": 0,
            "unchanged": 3
          }
        }
      }
    }
  }
}
```
