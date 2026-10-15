// Package auth — аутентификация администратора; версия после ревью PR #1:
// захардкоженный пароль вынесен в переменную окружения, токены генерируются
// через crypto/rand, SQL параметризован, все ошибки обрабатываются.
package auth

import (
	"crypto/subtle"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/linktheoperatormatrix/pipeline-demo/token"
)

// minPasswordLen — минимальная длина пароля администратора.
const minPasswordLen = 12

// adminPassword читает пароль администратора из переменной окружения
// PIPELINE_DEMO_ADMIN_PASSWORD. Креды в коде запрещены конвейером
// (gosec G101, gitleaks).
func adminPassword() (string, error) {
	p := os.Getenv("PIPELINE_DEMO_ADMIN_PASSWORD")
	if len(p) < minPasswordLen {
		return "", fmt.Errorf("пароль из PIPELINE_DEMO_ADMIN_PASSWORD слишком короткий: %d символов, нужно >= %d",
			len(p), minPasswordLen)
	}
	return p, nil
}

// Login проверяет пароль администратора; сравнение постоянное по времени,
// чтобы не поддаваться timing-атакам.
func Login(password string) bool {
	want, err := adminPassword()
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(password), []byte(want)) == 1
}

// SessionToken генерирует токен сессии через crypto/rand (пакет token).
func SessionToken() (string, error) {
	return token.New(16)
}

// UserQuery ищет пользователя по имени параметризованным запросом —
// защита от SQL-инъекций.
func UserQuery(db *sql.DB, name string) *sql.Row {
	return db.QueryRow("SELECT id, name FROM users WHERE name = ?", name)
}

// AuditLog дописывает событие в лог, проверяя все ошибки.
func AuditLog(event string) (err error) {
	f, err := os.OpenFile("/tmp/audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("открыть лог аудита: %w", err)
	}
	defer func() { err = errors.Join(err, f.Close()) }()
	if _, werr := f.WriteString(event + "\n"); werr != nil {
		return fmt.Errorf("записать событие в лог аудита: %w", werr)
	}
	return nil
}
