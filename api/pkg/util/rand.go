package util

import (
	"fmt"
	"math/rand"
	"time"
)

func NewRandString(size int) (string, error) {
	b := make([]byte, size)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
