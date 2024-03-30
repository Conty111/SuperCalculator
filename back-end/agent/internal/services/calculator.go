package services

import (
	"fmt"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/agent_errors"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

type CalculatorService struct {
	Locker          *sync.RWMutex
	AddTime         time.Duration
	MultiplyTime    time.Duration
	DivisionTime    time.Duration
	SubtractionTime time.Duration
	chars           []rune
}

func NewCalculatorService() *CalculatorService {
	return &CalculatorService{
		Locker: &sync.RWMutex{},
		chars:  []rune{'+', '-', '*', '/', '(', ')', '^', '.'},
	}
}

func (es *CalculatorService) Execute(t *models.Task) *models.Result {
	var r models.Result

	r.Task = *t

	expression, err := es.ValidateExpression(t.Expression)
	if err != nil {
		r.Error = err.Error()
		return &r
	}
	resNum, err := es.Calculate(expression)
	if err != nil {
		r.Error = err.Error()
		return &r
	}

	r.Value = resNum
	return &r
}

func (es *CalculatorService) SetOperationDuration(settings *models.DurationSettings) {
	es.Locker.RLock()
	defer es.Locker.RUnlock()
	es.DivisionTime = time.Millisecond * time.Duration(settings.DivisionDuration)
	es.MultiplyTime = time.Millisecond * time.Duration(settings.MultiplyDuration)
	es.SubtractionTime = time.Millisecond * time.Duration(settings.SubtractDuration)
	es.AddTime = time.Millisecond * time.Duration(settings.AddDuration)
}

// Calculate Calculating expression string
func (es *CalculatorService) Calculate(expression string) (float64, error) {
	rpn, err := es.infixToRPN(strings.TrimSpace(expression))
	if err != nil {
		return 0, err
	}
	res, err := es.evaluateRPN(rpn)
	if err != nil {
		return 0, err
	}

	return res, nil
}

// ValidateExpression Check if expression is valid and reformat it
func (es *CalculatorService) ValidateExpression(expression string) (string, error) {
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
			if i == 0 ||
				i == len(expression)-1 ||
				!unicode.IsDigit(rune(expression[i-1])) ||
				!unicode.IsDigit(rune(expression[i+1])) {

				return "", agent_errors.ErrInvalidFloat
			}
		} else if !unicode.IsDigit(char) && !slices.Contains(es.chars, char) {
			// Проверяем на неправильные вещественные числа (например, "2. * 3")
			return "", agent_errors.ErrInvalidChar(char)
		}
	}
	return expression, nil
}

func isOperator(ch string) bool {
	return ch == "+" || ch == "-" || ch == "*" || ch == "/"
}

func getPrecedence(operator string) int {
	switch operator {
	case "*", "/":
		return 1
	default:
		return 2
	}
}

func (es *CalculatorService) infixToRPN(infix string) ([]string, error) {
	var stack []string
	var nums []float64
	var operators []string
	var rpn []string

	// Функция для добавления оператора в результат или в стек
	pushOperator := func(operator string) {
		for len(stack) > 0 && getPrecedence(stack[len(stack)-1]) <= getPrecedence(operator) {
			operators = append(operators, stack[len(stack)-1])
			rpn = append(rpn, stack[len(stack)-1])
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, operator)
	}

	tokens := make([]string, 0)
	currentToken := ""

	// Разбиваем строку на токены
	for _, char := range infix {
		if unicode.IsDigit(char) || char == '.' || (char == '-' && (len(currentToken) == 0 || unicode.IsSpace(rune(currentToken[len(currentToken)-1])))) {
			currentToken += string(char)
		} else if isOperator(string(char)) || char == '(' || char == ')' {
			if currentToken != "" {
				tokens = append(tokens, currentToken)
				currentToken = ""
			}
			tokens = append(tokens, string(char))
		}
	}
	if currentToken != "" {
		tokens = append(tokens, currentToken)
	}

	// Преобразуем токены в обратную польскую запись
	for i := 0; i < len(tokens); i++ {
		if num, err := strconv.ParseFloat(tokens[i], 64); err == nil {
			// Если токен - число, добавляем его в результат
			//if indexToInsert > -1 {
			//	for j := indexToInsert; j < len(nums); j++ {
			//		nums[j], num = num, nums[j]
			//	}
			//}
			nums = append(nums, num)
			rpn = append(rpn, tokens[i])
		} else if isOperator(tokens[i]) {
			// Если токен - оператор, добавляем его в результат или в стек
			pushOperator(tokens[i])
		} else if tokens[i] == "(" {
			// Если токен - '(', запускаем рекурсию
			var num float64
			var isFinded bool
			for j := i + 1; j < len(tokens); j++ {
				if tokens[j] == ")" {
					isFinded = true
					num, err = es.Calculate(strings.Join(tokens[i+1:j], ""))
					if err != nil {
						return rpn, err
					}
					i = j
					break
				}
			}
			if !isFinded {
				return rpn, agent_errors.ErrParenthesisNotValid
			}
			nums = append(nums, num)
			rpn = append(rpn, fmt.Sprintf("%f", num))
		}
	}

	// Добавляем оставшиеся операторы из стека в результат
	for len(stack) > 0 {
		operators = append(operators, stack[len(stack)-1])
		rpn = append(rpn, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	if len(nums)-len(operators) != 1 {
		return rpn, agent_errors.ErrInvalidExpression
	}

	return rpn, nil
}

func (es *CalculatorService) evaluateRPN(rpn []string) (float64, error) {
	var stack []float64

	for _, token := range rpn {
		if num, err := strconv.ParseFloat(token, 64); err != nil {
			operand2 := stack[len(stack)-1]
			operand1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			var result float64
			switch token {
			case "+":
				time.Sleep(es.AddTime)
				result = operand1 + operand2
			case "-":
				time.Sleep(es.SubtractionTime)
				result = operand1 - operand2
			case "*":
				time.Sleep(es.MultiplyTime)
				result = operand1 * operand2
			case "/":
				time.Sleep(es.DivisionTime)
				if operand2 == 0 {
					return 0, agent_errors.ErrDivisionByZero
				}
				result = operand1 / operand2
			}
			stack = append(stack, result)
		} else {
			stack = append(stack, num)
		}
	}

	return stack[0], nil
}
