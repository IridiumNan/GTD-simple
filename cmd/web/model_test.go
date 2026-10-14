package web

import (
	"encoding/json"
	"testing"
)

func TestJSONUnmarshal(t *testing.T) {
	data := `

{
    "data": {
        "stuff": {
            "id": "",
            "raw": "this is raw content for this stuff",
            "is_del": false
        }
    }
}
	`
	var req RequestCapture

	err := json.Unmarshal([]byte(data), &req)
	if err != nil {
		t.Errorf("error while unmarshal test josn string, err: %s", err.Error())
	}
}
