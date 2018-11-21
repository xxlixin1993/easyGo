package random

import (
	"math/rand"
	"strings"
	"time"
)

type (
	Random struct {
	}
)

// Charsets
const (
	Uppercase    string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Lowercase           = "abcdefghijklmnopqrstuvwxyz"
	Alphabetic          = Uppercase + Lowercase
	Numeric             = "0123456789"
	Alphanumeric        = Alphabetic + Numeric
	Symbols             = "`" + `~!@#$%^&*()-_+={}[]|\;:"<>,./?`
	Hex                 = Numeric + "abcdef"
)

var (
	global = New()
)

func New() *Random {
	rand.Seed(time.Now().UnixNano())
	return new(Random)
}

func (r *Random) String(length uint8, charsets ...string) string {
	charset := strings.Join(charsets, "")
	if charset == "" {
		charset = Alphanumeric
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Int63()%int64(len(charset))]
	}
	return string(b)
}

func RandomString(length uint8, charsets ...string) string {
	return global.String(length, charsets...)
}

func RandUInt32(minNum, maxNum uint32) uint32 {
	min := int32(minNum)
	max := int32(maxNum)
	if min >= max || min == 0 || max == 0 {
		return uint32(max)
	}
	return uint32(rand.Int31n(max-min) + min)
}