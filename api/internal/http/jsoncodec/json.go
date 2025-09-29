package jsoncodec

import "encoding/json"

func MustMarshal(v any) json.RawMessage {
   b, err := json.Marshal(v)
   if err != nil {
       return json.RawMessage{}
   }
   return json.RawMessage(b)
}
