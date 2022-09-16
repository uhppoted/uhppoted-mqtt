{{- with .acl_upload -}}
### `{{.command}}`

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

{{template "request"  . -}}
{{template "response" . }}

Example:
```
{
  "message": {
    "request": {
{{- template "request-preamble"}}
    "request": {
      "url": "file:///var/uhppoted/uploaded/ACL.tar.gz",
    }
  }
}

{
  "message": {
    "reply": {
{{- template "response-preamble"}}
      "method": "acl:upload",
      "response": {
        "uploaded": "file:///var/uhppoted/uploaded/uhppoted.tar.gz"
      }
    }
  }
}
```
{{end -}}
