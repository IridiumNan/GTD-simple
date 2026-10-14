package models

import (
	"strings"
	"time"

	"github.com/IridiumNan/GTD-simple/internal/utils"
)

type Stuff struct {
	ID    string `json:"id"`
	Raw   string `json:"raw"`
	IsDel bool   `json:"is_del"`
}

// genIDByNow generate new string based on current time
func genIDByNow() string {
	return utils.HashByTimestamp(time.Now().UnixNano())
}

// NewStuff create a new stuff with raw content
func NewStuff(raw string) *Stuff {
	return &Stuff{
		ID:    genIDByNow(),
		Raw:   strings.TrimSpace(raw),
		IsDel: false,
	}
}
