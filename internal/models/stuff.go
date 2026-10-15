package models

import (
	"encoding/json"
	"fmt"
	"log/slog"
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

type StuffStat struct {
	Count int   `json:"count"`
	Size  int64 `json:"size"`
}

func (ss *StuffStat) String() string {
	return fmt.Sprintf("%d stuff with size: %s", ss.Count, utils.FormatBytes(ss.Size))
}

func (ss *StuffStat) JSON() []byte {
	data, err := json.Marshal(ss)
	if err != nil {
		slog.Error("error while unmarshal stuff stat into json byte data", "err", err, "stat", ss)
		return nil
	}

	return data
}
