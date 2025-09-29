package utils

import (
	"encoding/json"
)

func MustMarshal(v interface{}) json.RawMessage {
    b, _ := json.Marshal(v)
    return json.RawMessage(b)
}