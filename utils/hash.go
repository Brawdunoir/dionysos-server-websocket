package utils

import (
	"crypto/sha1"
	"fmt"
)

// Generate a hash from a string
func GenerateStringHash(s string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}
