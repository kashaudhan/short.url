package utils

import (
	"errors"
	"math"
	"strings"
)

const (
	chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func Base62Encode(num uint64) string {
	length := len(chars)

	var encodedBuilder strings.Builder

	encodedBuilder.Grow(10)

	for ; num > 0; num = num / uint64(length) {
		encodedBuilder.WriteByte(chars[(num % uint64(length))])
	}

	return encodedBuilder.String()
}

func Base62Decode(encodedString string) (uint64, error) {
	var num uint64
	length := len(encodedString)

	for i, char := range encodedString {
		position := strings.IndexRune(chars, char)

		if position == -1 {
			return uint64(position), errors.New("invalid character")
		}

		num += uint64(position) * uint64(math.Pow(float64(length), float64(i)))
	}

	return num, nil
}