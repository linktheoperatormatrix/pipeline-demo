// Package calc — простые арифметические операции для демо-конвейера.
package calc

import "errors"

// ErrDivisionByZero возвращает Divide при делении на ноль.
var ErrDivisionByZero = errors.New("division by zero")

// Add возвращает сумму a и b.
func Add(a, b float64) float64 { return a + b }

// Subtract возвращает разность a и b.
func Subtract(a, b float64) float64 { return a - b }

// Divide делит a на b и возвращает ErrDivisionByZero, если b равен нулю.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return a / b, nil
}
