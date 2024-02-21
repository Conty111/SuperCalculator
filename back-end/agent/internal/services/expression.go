package services

import (
	"encoding/json"
	"fmt"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/agent_errors"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

type ExpressionService struct {
	Locker          *sync.RWMutex
	AddTime         time.Duration
	MultiplyTime    time.Duration
	DivisionTime    time.Duration
	SubtractionTime time.Duration
	chars           []rune
}

func NewExpressionService() *ExpressionService {
	return &ExpressionService{
		Locker: &sync.RWMutex{},
		chars:  []rune{'+', '-', '*', '/', '(', ')', '^', '.'},
	}
}

func (es *ExpressionService) Proccess(msg *sarama.ConsumerMessage) *models.Result {
	var t models.Task
	var r models.Result
	if err := json.Unmarshal(msg.Value, &t); err != nil {
		log.Error().Msg("Error while parsing json")
		return nil
	}
	r.Task = t
	log.Debug().Any("task", t).Str("key", string(msg.Key)).Msg("parsed to json")
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

func (es *ExpressionService) SetOperationDuration(settings *models.DurationSettings) {
	es.Locker.RLock()
	defer es.Locker.RUnlock()
	es.DivisionTime = time.Millisecond * time.Duration(settings.DivisionDuration)
	es.MultiplyTime = time.Millisecond * time.Duration(settings.MultiplyDuration)
	es.SubtractionTime = time.Millisecond * time.Duration(settings.SubtractDuration)
	es.AddTime = time.Millisecond * time.Duration(settings.AddDuration)
}

func (es *ExpressionService) Calculate(expression string) (float64, error) {
	exp := es.infixToRPN(strings.TrimSpace(expression))
	res, err := es.evaluateRPN(exp)
	if err != nil {
		return 0, err
	}

	return res, nil
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
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

func (es *ExpressionService) infixToRPN(infix string) string {
	var result string
	var stack []string

	// Функция для добавления оператора в результат или в стек
	pushOperator := func(operator string) {
		for len(stack) > 0 && isOperator(stack[len(stack)-1]) && getPrecedence(stack[len(stack)-1]) >= getPrecedence(operator) {
			result += stack[len(stack)-1] + " "
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
		} else if !unicode.IsSpace(char) {
			// Пропускаем пробелы и обрабатываем только числа и операторы
			fmt.Printf("Unsupported character: %c\n", char)
		}
	}
	if currentToken != "" {
		tokens = append(tokens, currentToken)
	}

	// Преобразуем токены в обратную польскую запись
	for _, token := range tokens {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			// Если токен - число, добавляем его в результат
			result += fmt.Sprintf("%f ", num)
		} else if isOperator(token) {
			// Если токен - оператор, добавляем его в результат или в стек
			pushOperator(token)
		} else if token == "(" {
			// Если токен - '(', добавляем его в стек
			stack = append(stack, token)
		} else if token == ")" {
			// Если токен - ')', перемещаем операторы из стека в результат до '('
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				result += stack[len(stack)-1] + " "
				stack = stack[:len(stack)-1]
			}
			// Убираем '(' из стека
			stack = stack[:len(stack)-1]
		}
	}

	// Добавляем оставшиеся операторы из стека в результат
	for len(stack) > 0 {
		result += stack[len(stack)-1] + " "
		stack = stack[:len(stack)-1]
	}

	return strings.TrimSpace(result)
}

func (es *ExpressionService) evaluateRPN(rpn string) (float64, error) {
	var stack []float64
	tokens := strings.Fields(rpn)

	for _, token := range tokens {
		if token == "" {
			continue // Игнорировать пустые строки
		}

		if num, err := strconv.ParseFloat(token, 64); err == nil {
			// Если токен - число, поместить его в стек
			stack = append(stack, num)
		} else if isOperator(token) {
			// Если токен - оператор, выполнить операцию с двумя верхними элементами стека
			if len(stack) < 2 {
				return 0, fmt.Errorf("insufficient operands for operator %s", token)
			}
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
					return 0, fmt.Errorf("division by zero")
				}
				result = operand1 / operand2
			}

			stack = append(stack, result)
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid expression")
	}

	return stack[0], nil
}
