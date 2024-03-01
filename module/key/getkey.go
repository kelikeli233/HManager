package key

import (
	"encoding/base64"
	"math/rand"
)

func RandomKey(length ...int) (string, error) {
	var key []byte
	if len(length) == 0 {
		key = make([]byte, 64)
	} else {
		key = make([]byte, length[0])
	}
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}
