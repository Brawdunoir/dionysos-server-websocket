package objects

import (
	"crypto/sha1"
	"fmt"
)

// Generate a hash from a string
func generateStringHash(s string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}
