package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"strings"
)

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

//data -> byte
func ToBytes(i interface{}) []byte {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	HandleErr(encoder.Encode(i))
	return aBuffer.Bytes()
}

//byte -> data
func FromBytes(i interface{}, data []byte) {
	encoder := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(encoder.Decode(i))
}

//data -> hash
func Hash(a interface{}) string {
	str := fmt.Sprintf("%v", a)
	hash := sha256.Sum256([]byte(str))
	return fmt.Sprintf("%x", hash)
}

func IPSplitter(s string, sep string, i int) string {
	result := strings.Split(s, sep)
	if len(result)-1 < i {
		return ""
	}
	return result[i]
}
