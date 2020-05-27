package acl

import (
	"path/filepath"
)

var DEFAULT_KEYFILE = filepath.Join(workdir(), "acl", "keys", "uhppoted")
var DEFAULT_CREDENTIALS = filepath.Join(workdir(), ".aws", "credentials")
var DEFAULT_REGION = "us-east-1"
