package agent_errors

import (
	"errors"
	"fmt"
)

var (
	ErrParenthesisNotValid = errors.New("несоответствие количества открытых и закрытых скобок")
	ErrTooFewOperands      = errors.New("недостаточно чисел")
	ErrTooManyOperands     = errors.New("лишние числа")
)

func ErrInvalidChar(char rune) error {
	return fmt.Errorf("недопустимый символ: %s", char)
}

func ErrOperandParsing(operandString string, err error) error {
	return fmt.Errorf("ошибка при парсинге операнда %s: %v", operandString, err)
}
