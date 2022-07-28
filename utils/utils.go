// Package utils contains functions to be used across application

package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
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
//FromBytes takes an interface and data and the will encode the data to the interface
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

func Splitter(s string, sep string, i int) string {
	result := strings.Split(s, sep)
	if len(result)-1 < i {
		return ""
	}
	return result[i]
}

//json -> byte
func ToJSON(i interface{}) []byte {
	r, err := json.Marshal(i)
	HandleErr(err)
	return r
}
