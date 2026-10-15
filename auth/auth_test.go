package auth

import (
	"testing"
)

func TestAdminPasswordTooShort(t *testing.T) {
	t.Setenv("PIPELINE_DEMO_ADMIN_PASSWORD", "short")
	if _, err := adminPassword(); err == nil {
		t.Fatal("короткий пароль должен отклоняться")
	}
}

func TestAdminPasswordOK(t *testing.T) {
	t.Setenv("PIPELINE_DEMO_ADMIN_PASSWORD", "correct-horse-battery")
	if _, err := adminPassword(); err != nil {
		t.Fatalf("корректный пароль отклонён: %v", err)
	}
}

func TestLogin(t *testing.T) {
	t.Setenv("PIPELINE_DEMO_ADMIN_PASSWORD", "correct-horse-battery")
	if !Login("correct-horse-battery") {
		t.Fatal("верный пароль отклонён")
	}
	if Login("wrong-password-123") {
		t.Fatal("неверный пароль принят")
	}
}

func TestLoginWithoutPassword(t *testing.T) {
	t.Setenv("PIPELINE_DEMO_ADMIN_PASSWORD", "")
	if Login("whatever-long-enough") {
		t.Fatal("без настроенного пароля вход должен отклоняться")
	}
}

func TestSessionToken(t *testing.T) {
	a, err := SessionToken()
	if err != nil {
		t.Fatalf("SessionToken: %v", err)
	}
	b, err := SessionToken()
	if err != nil {
		t.Fatalf("SessionToken: %v", err)
	}
	if len(a) != 32 {
		t.Fatalf("токен должен быть 32 hex-символа, получен %q", a)
	}
	if a == b {
		t.Fatalf("токены должны быть уникальными: %q == %q", a, b)
	}
}
