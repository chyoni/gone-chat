package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
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

func ToUintFromString(aString string) uint {
	aStringAsUint, err := strconv.ParseUint(aString, 10, 64)
	if err != nil {
		HandleError(err)
	}
	return uint(aStringAsUint)
}

func ConnectAws() *session.Session {
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	myRegion := os.Getenv("AWS_REGION")

	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(myRegion),
			Credentials: credentials.NewStaticCredentials(
				accessKeyID,
				secretAccessKey,
				"",
			),
		})
	if err != nil {
		HandleError(err)
	}
	return sess
}
