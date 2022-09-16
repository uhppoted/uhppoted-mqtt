{{- with .acl_compare -}}
### `{{.command}}`

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

{{template "request"  . -}}
{{template "response" . }}

Example:
```
{
  "message": {
    "request": {
{{- template "request-preamble"}}
    "request": {
      "acl": "file:///var/uhppoted/ACL.tar.gz",
      "report": "file:///var/uhppoted/report.tar.gz"
    }
  }
}

{
  "message": {
    "reply": {
{{- template "response-preamble"}}
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
{{end -}}
