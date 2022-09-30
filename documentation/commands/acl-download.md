### `acl-download`

Updates the access control permissions on a set of controllers from a file retrieved from a URL. The URL
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

topic: <root>/<requests>/acl/acl:download

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "url": "URL",
            "mime-type": "IANA mime-type",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
url          Source URL for the ACL file
mime-type    application/x-gzip for tar.gz files, application/zip for .zip files and text/tab-separated-values for TSV files
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "acl-download",
      "response": {
            "report": "list of record",
            "report.controller": "uint32",
            "report.controller.added": "uint32",
            "report.controller.deleted": "uint32",
            "report.controller.updated": "uint32",
            "report.controller.unchanged": "uint32",
            "report.controller.errors": "uint32",
            "report.controller.failed": "uint32",
            "warnings": "array of string",
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
report       list of the changes made to each controller
report.controller controller ID
report.controller.added number of cards added to controller
report.controller.deleted number of cards deleted from controller
report.controller.updated number of cards updated on controller
report.controller.unchanged number of cards unchanged on controller
report.controller.errors number of errors for controller
report.controller.failed number of cards that could not be transferred to controller
warnings     list of warning messages while transferring ACL to controller
```


Example:
```
topic: uhppoted/gateway/requests/acl/acl:download

{
  "message": {
    "request": {
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
    "request": {
      "url": "file:///var/uhppoted/ACL.tar.gz",
      "mime-type": "application/x-gzip"
    }
  }
}

{
  "message": {
    "reply": {
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "method": "acl:download",
      "response": {
        "report": {
          "201020304": {
            "added": 3,
            "deleted": 0,
            "errors": 0,
            "failed": 0,
            "unchanged": 0,
            "updated": 0
          },
          "303986753": {
            "added": 1,
            "deleted": 0,
            "errors": 2,
            "failed": 0,
            "unchanged": 0,
            "updated": 0
          },
          "405419896": {
            "added": 0,
            "deleted": 0,
            "errors": 0,
            "failed": 0,
            "unchanged": 0,
            "updated": 3
          }
        },
        "warnings": [
          "303986753: Time profile 29 is not defined for 303986753"
        ]
      }
    }
  }
}
```
