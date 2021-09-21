package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToBytes(v interface{}) []byte {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	err := encoder.Encode(v)
	HandleError(err)
	return aBuffer.Bytes()
}

func FromBytes(v interface{}, data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(v)
	HandleError(err)
}

func ToHexStringHash(data []byte) string {
	hash := sha256.Sum256(data)
	hexHash := fmt.Sprintf("%x", hash)
	return hexHash
}
