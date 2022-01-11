package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
)

func same(a []byte, b string) bool {
	return bytes.Equal(a, []byte(b))
}

// Generate a hash from a string
func generateStringHash(s string) string {
	return fmt.Sprint(sha1.Sum([]byte(s)))
}
