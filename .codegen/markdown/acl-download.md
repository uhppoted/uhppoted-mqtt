{{- with .acl_download -}}
### `{{.command}}`

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

{{template "request"  . -}}
{{template "response" . }}

Example:
```
topic: uhppoted/gateway/requests/{{ .request.topic }}

{
  "message": {
    "request": {
{{- template "request-preamble"}}
    "request": {
      "url": "file:///var/uhppoted/ACL.tar.gz",
      "mime-type": "application/x-gzip"
    }
  }
}

{
  "message": {
    "reply": {
{{- template "response-preamble"}}
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
{{end -}}
