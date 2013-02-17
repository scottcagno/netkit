package netkit

import (
	"encoding/json"
	"encoding/hex"
)

// map type
type M map[string]interface{}

// hex encoder
func EncodeHex(s string) string {
	return hex.EncodeToString([]byte(s))
}

// hex decoder
func DecodeHex(s string) string {
	dat, _ := hex.DecodeString(s)
	return string(dat)
}

// json encoder
func EncodeJSON(v interface{}) string {
	dat, _ := json.Marshal(v)
	return string(dat)
}

// json encoder
func EncodePrettyJSON(v interface{}) string {
	dat, _ := json.MarshalIndent(v, "", "    ")
	return string(dat)
}

// json decoder
func DecodeJSON(s string, v interface{}) {
	json.Unmarshal([]byte(s), &v)
}

// marshaler
func Marshal(v1, v2 interface{}) {
	dat, _ := json.Marshal(v1)
	json.Unmarshal(dat, &v2)
}