package token

import "testing"

func TestNewLength(t *testing.T) {
	t.Parallel()
	got, err := New(16)
	if err != nil {
		t.Fatalf("New(16) unexpected error: %v", err)
	}
	if len(got) != 32 {
		t.Fatalf("New(16) length = %d, want 32", len(got))
	}
}

func TestNewUnique(t *testing.T) {
	t.Parallel()
	a, err := New(8)
	if err != nil {
		t.Fatalf("New(8) unexpected error: %v", err)
	}
	b, err := New(8)
	if err != nil {
		t.Fatalf("New(8) unexpected error: %v", err)
	}
	if a == b {
		t.Fatalf("токены должны быть уникальными: %q == %q", a, b)
	}
}
