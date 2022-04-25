package utils

import (
	"Assignment/models"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"time"
	"unsafe"
)

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits

)

var (
	alphabet = []byte(`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890`)
	src      = rand.NewSource(time.Now().UnixNano())
)

func EncodeJSONBody(resp http.ResponseWriter, statusCode int, data interface{}) {
	resp.WriteHeader(statusCode)
	err := json.NewEncoder(resp).Encode(data)
	if err != nil {
		logrus.Errorf("Error encoding response %v", err)
	}
}

func EncodeJSON200Body(resp http.ResponseWriter, data interface{}) {
	err := json.NewEncoder(resp).Encode(data)
	if err != nil {
		logrus.Errorf("Error encoding response %v", err)
	}
}

func CreateInviteCode(cache map[string]interface{}) (string, error) {
	maxTry := 5
	var inviteUniqueCode string

	for maxTry > 0 {
		inviteUniqueCode = generateRandomCode(models.InviteCodeLength)

		var isAlreadyExists bool
		if _, ok := cache[inviteUniqueCode]; ok {
			isAlreadyExists = true
		}

		if len(inviteUniqueCode) == models.InviteCodeLength && !isAlreadyExists {
			break
		}

		maxTry--
		if maxTry == 0 {
			return "", errors.New("max limit exceeded")
		}
	}

	return inviteUniqueCode, nil
}

func generateRandomCode(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(alphabet) {
			b[i] = alphabet[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
