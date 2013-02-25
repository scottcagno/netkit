package netkit

import (
	"encoding/base64"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
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

// random id generator
func RandomId() string {
	e := make([]byte, 32)
	rand.Read(e)
	b := make([]byte, base64.URLEncoding.EncodedLen(len(e)))
	base64.URLEncoding.Encode(b, e)
	return string(b)
}

// url encode a string
func URLEncode(s string) string {
	encoder := base64.URLEncoding
	encoded := make([]byte, encoder.EncodedLen(len([]byte(s))))
	encoder.Encode(encoded, []byte(s))
	return string(encoded)
}

// url decode a string
func URLDecode(s string) string {
	encoder := base64.URLEncoding
	decoded := make([]byte, encoder.EncodedLen(len([]byte(s))))
	_, err := encoder.Decode(decoded, []byte(s))
	if err != nil {
		return fmt.Sprintln(err)
	}
	return string(decoded)
}
