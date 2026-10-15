package calc

import (
	"errors"
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	t.Parallel()
	if got := Add(2, 3); got != 5 {
		t.Fatalf("Add(2, 3) = %v, want 5", got)
	}
}

func TestSubtract(t *testing.T) {
	t.Parallel()
	if got := Subtract(5, 3); got != 2 {
		t.Fatalf("Subtract(5, 3) = %v, want 2", got)
	}
}

func TestDivide(t *testing.T) {
	t.Parallel()
	got, err := Divide(10, 4)
	if err != nil {
		t.Fatalf("Divide(10, 4) unexpected error: %v", err)
	}
	if math.Abs(got-2.5) > 1e-9 {
		t.Fatalf("Divide(10, 4) = %v, want 2.5", got)
	}
}

func TestDivideByZero(t *testing.T) {
	t.Parallel()
	_, err := Divide(1, 0)
	if !errors.Is(err, ErrDivisionByZero) {
		t.Fatalf("Divide(1, 0) error = %v, want ErrDivisionByZero", err)
	}
}
