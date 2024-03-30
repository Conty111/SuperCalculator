package agent_errors

import (
	"errors"
	"fmt"
)

var (
	ErrParenthesisNotValid = errors.New("несоответствие количества открытых и закрытых скобок")
	ErrInvalidExpression   = errors.New("неправильное выражение (кол-во операндов и операций неправильное)")
	ErrInvalidFloat        = errors.New("неправильное вещественное число")
	ErrInvalidMessage      = errors.New("неверный формат сообщения (должен быть json)")
	ErrDivisionByZero      = errors.New("деление на 0, ты чего")
)

const NotAllowedChar = "недопустимый символ"

func ErrInvalidChar(char rune) error {
	return fmt.Errorf("%s: %v", NotAllowedChar, char)
}

func ErrOperandParsing(operandString string, err error) error {
	return fmt.Errorf("ошибка при парсинге операнда %s: %v", operandString, err)
}
