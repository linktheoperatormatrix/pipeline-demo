// Package auth — ДЕМО-версия с типовыми уязвимостями:
// конвейер (хуки локально, gosec/golangci-lint в CI) должен их поймать.
package auth

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
)

// Захардкоженный пароль — gosec G101, потенциально gitleaks.
const adminPassword = "hunter22"

// Login проверяет пароль администратора.
func Login(password string) bool {
	return password == adminPassword
}

// SessionToken генерирует токен сессии через math/rand —
// слабая криптография, gosec G404.
func SessionToken() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// UserQuery ищет пользователя по имени конкатенацией SQL —
// инъекция, gosec G201.
func UserQuery(db *sql.DB, name string) *sql.Row {
	q := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", name)
	return db.QueryRow(q)
}

// AuditLog дописывает событие в лог, игнорируя ошибки — errcheck.
func AuditLog(event string) {
	f, _ := os.OpenFile("/tmp/audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString(event + "\n")
}
