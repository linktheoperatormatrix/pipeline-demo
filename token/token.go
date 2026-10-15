// Package token генерирует криптографически стойкие токены сессий.
package token

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// New возвращает случайный токен из n байт в hex-кодировке (2n символов).
func New(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}
	return hex.EncodeToString(b), nil
}
