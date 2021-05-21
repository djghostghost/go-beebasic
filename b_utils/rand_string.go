package b_utils

import (
	"fmt"
	"math/rand"
)

func RandString(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}
