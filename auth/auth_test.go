package auth

import "testing"

// Требование безопасности: пароль администратора не короче 12 символов.
// Этот тест упадёт — и хук pre-push, и CI обязаны это поймать.
func TestAdminPasswordStrength(t *testing.T) {
	if len(adminPassword) < 12 {
		t.Fatalf("пароль администратора слишком короткий: %d символов, нужно >= 12",
			len(adminPassword))
	}
}

func TestLogin(t *testing.T) {
	if !Login("hunter22") {
		t.Fatal("верный пароль отклонён")
	}
	if Login("wrong") {
		t.Fatal("неверный пароль принят")
	}
}
