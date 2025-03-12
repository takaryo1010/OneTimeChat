package controller

import (
	"math/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateSessionID() (string, error) {

	data := fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Int())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:]), nil
}
