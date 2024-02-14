package services

import (
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/agent_errors"
	"math"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

type ExpressionService struct {
	Locker             *sync.Mutex
	AddTime            time.Duration
	MultiplyTime       time.Duration
	DivisionTime       time.Duration
	SubtractionTime    time.Duration
	ExponentiationTime time.Duration
	chars              []rune
}

func NewExpressionService() *ExpressionService {
	return &ExpressionService{
		Locker: &sync.Mutex{},
		chars:  []rune{'+', '-', '*', '/', '(', ')', '^', '.'},
	}
}

func (es *ExpressionService) SetOperationDuration(operation rune, t time.Duration) error {
	es.Locker.Lock()
	defer es.Locker.Unlock()
	switch operation {
	case '-':
		es.SubtractionTime = t
	case '+':
		es.AddTime = t
	case '/':
		es.DivisionTime = t
	case '*':
		es.MultiplyTime = t
	case '^':
		es.ExponentiationTime = t
	default:
		return agent_errors.ErrInvalidChar(operation)
	}
	return nil
}

func (es *ExpressionService) Calculate(expression string) (float64, error) {
	exp, err := es.ValidateExpression(expression)
	if err != nil {
		return 0, err
	}

	operands, operators, err := es.ParseToInfix(exp)
	if err != nil {
		return 0, err
	}

	if len(operators) < len(operands)-1 {
		return 0, agent_errors.ErrTooManyOperands
	} else if len(operators) > len(operands)-1 {
		return 0, agent_errors.ErrTooFewOperands
	}

	result := es.CalculateInfix(operands, operators)
	return result, nil
}

// ValidateExpression Check if expression is valid and reformat it
func (es *ExpressionService) ValidateExpression(expression string) (string, error) {
	// Проверяем наличие соответствия открытых и закрытых скобок
	expression = strings.ReplaceAll(expression, " ", "")
	openCount, closeCount := 0, 0
	for _, char := range expression {
		switch char {
		case '(':
			openCount++
		case ')':
			closeCount++
		}
	}

	if openCount != closeCount {
		return "", agent_errors.ErrParenthesisNotValid
	}

	// Проверяем корректность символов в выражении
	for i, char := range expression {
		if char == '.' {
			// Проверяем, что точка не является первым или последним символом, и перед и после точки есть цифры
			if i == 0 || i == len(expression)-1 || !unicode.IsDigit(rune(expression[i-1])) || !unicode.IsDigit(rune(expression[i+1])) {
				return "", agent_errors.ErrInvalidFloat
			}
		} else if !unicode.IsDigit(char) && !slices.Contains(es.chars, char) {
			// Проверяем на неправильные вещественные числа (например, "2. * 3")
			return "", agent_errors.ErrInvalidChar(char)
		}
	}

	return expression, nil
}

// Returns operands and operators from expression in preorder
func (es *ExpressionService) ParseToInfix(expression string) ([]float64, []string, error) {
	// Разбираем выражение на операторы и операнды
	operators := make([]string, 0)
	operands := make([]float64, 0)

	i := 0
	for i < len(expression) {
		switch expression[i] {
		case '+', '-', '*', '/', '^':
			operators = append(operators, string(expression[i]))
			i++
		case '(':
			// Ищем соответствующую закрывающую скобку
			openCount := 1
			closeIndex := i + 1
			for closeIndex < len(expression) && openCount > 0 {
				if expression[closeIndex] == '(' {
					openCount++
				} else if expression[closeIndex] == ')' {
					openCount--
				}
				closeIndex++
			}

			// Рекурсивно вычисляем значение внутри скобок
			subExpression := expression[i+1 : closeIndex-1]
			subResult, err := es.Calculate(subExpression)
			if err != nil {
				return nil, nil, err
			}
			operands = append(operands, subResult)
			i = closeIndex
		default:
			// Читаем операнд
			var operandString string
			for i < len(expression) && (expression[i] == '.' || unicode.IsDigit(rune(expression[i]))) {
				operandString += string(expression[i])
				i++
			}

			operand, err := strconv.ParseFloat(operandString, 64)
			if err != nil {
				return nil, nil, agent_errors.ErrOperandParsing(operandString, err)
			}
			operands = append(operands, operand)
		}
	}
	return operands, operators, nil
}

// Calculate preorder operands and operators
func (es *ExpressionService) CalculateInfix(operands []float64, operators []string) float64 {
	// Вычисляем выражение с учетом задержек для всех операций
	result := operands[0]

	for i := 0; i < len(operators); i++ {
		switch operators[i] {
		case "+":
			// Задержка для учета времени выполнения сложения
			time.Sleep(es.AddTime)
			result += operands[i+1]
		case "-":
			// Задержка для учета времени выполнения вычитания
			time.Sleep(es.SubtractionTime)
			result -= operands[i+1]
		case "*":
			// Задержка для учета времени выполнения умножения
			time.Sleep(es.MultiplyTime)
			result *= operands[i+1]
		case "/":
			// Задержка для учета времени выполнения деления
			time.Sleep(es.DivisionTime)
			result /= operands[i+1]
		case "^":
			// Задержка для учета времени выполнения возведения в степень
			time.Sleep(es.ExponentiationTime)
			result = math.Pow(result, operands[i+1])
		}
	}

	return result
}
