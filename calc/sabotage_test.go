package calc

import "testing"

// Демонстрация замка мержа: этот тест падает.
func TestSabotage(t *testing.T) {
	t.Fatal("намеренно красный тест: мерж должен быть заблокирован")
}
