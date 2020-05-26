package acl

import (
	"github.com/uhppoted/uhppote-core/types"
)

type Permission struct {
	Door      string     `json:"door"`
	StartDate types.Date `json:"start-date"`
	EndDate   types.Date `json:"end-date"`
}

type Permissions struct {
	CardNumber  uint32       `json:"card-number"`
	Permissions []Permission `json:"permissions"`
}

type Error struct {
	Code    int    `json:"error-code"`
	Message string `json:"message"`
}
