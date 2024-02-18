package services

import (
	"encoding/json"
	"fmt"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/agent_errors"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
	"math"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

type ExpressionService struct {
	Locker          *sync.Mutex
	AddTime         time.Duration
	MultiplyTime    time.Duration
	DivisionTime    time.Duration
	SubtractionTime time.Duration
	chars           []rune
}

func NewExpressionService() *ExpressionService {
	return &ExpressionService{
		Locker: &sync.Mutex{},
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
	//key, err := strconv.Atoi(string(msg.Key))
	//if err != nil {
	//	return nil, err
	//}
	//t.ID = uint(key)
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
	es.Locker.Lock()
	defer es.Locker.Unlock()
	es.DivisionTime = time.Millisecond * time.Duration(settings.DivisionDuration)
	es.MultiplyTime = time.Millisecond * time.Duration(settings.MultiplyDuration)
	es.SubtractionTime = time.Millisecond * time.Duration(settings.SubtractDuration)
	es.AddTime = time.Millisecond * time.Duration(settings.AddDuration)
}

func (es *ExpressionService) Calculate(expression string) (float64, error) {
	infixForm := es.ToInfix(expression)
	result, err := es.CalculateInfix(infixForm)
	if err != nil {
		return 0, err
	}
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

func (es *ExpressionService) ToInfix(expr string) []string {
	output := []string{}
	operators := map[string]int{"+": 1, "-": 1, "*": 2, "/": 2, "^": 3}

	tokens := tokenize(expr)
	stack := []string{}

	for _, token := range tokens {
		if isNumber(token) || isNegativeNumber(token) {
			output = append(output, token)
		} else if token == "(" {
			stack = append(stack, token)
		} else if token == ")" {
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1] // Remove "("
		} else if isOperator(token) {
			for len(stack) > 0 && operators[token] <= operators[stack[len(stack)-1]] {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		}
	}

	for len(stack) > 0 {
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output
}

func (es *ExpressionService) CalculateInfix(tokens []string) (float64, error) {
	stack := []float64{}

	for _, token := range tokens {
		if isNumber(token) || isNegativeNumber(token) {
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, err
			}
			stack = append(stack, num)
		} else if isOperator(token) {
			if len(stack) < 2 {
				return 0, fmt.Errorf("insufficient operands for operator: %s", token)
			}
			operand2 := stack[len(stack)-1]
			operand1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch token {
			case "+":
				time.Sleep(es.AddTime)
				stack = append(stack, operand1+operand2)
			case "-":
				time.Sleep(es.SubtractionTime)
				stack = append(stack, operand1-operand2)
			case "*":
				time.Sleep(es.MultiplyTime)
				stack = append(stack, operand1*operand2)
			case "/":
				time.Sleep(es.DivisionTime)
				if operand2 == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				stack = append(stack, operand1/operand2)
			case "^":
				time.Sleep(es.MultiplyTime * time.Duration(operand2))
				stack = append(stack, math.Pow(operand1, operand2))
			}
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid expression")
	}

	return stack[0], nil
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/" || s == "^"
}

func isNegativeNumber(s string) bool {
	return strings.HasPrefix(s, "-") && len(s) > 1 && isNumber(s[1:])
}

// returns slice of operators
func tokenize(expr string) []string {
	tokens := []string{}
	currentToken := ""

	for _, char := range expr {
		if char == ' ' {
			continue
		}

		if isOperator(string(char)) || char == '(' || char == ')' {
			if currentToken != "" {
				tokens = append(tokens, currentToken)
				currentToken = ""
			}
			tokens = append(tokens, string(char))
		} else {
			currentToken += string(char)
		}
	}

	if currentToken != "" {
		tokens = append(tokens, currentToken)
	}

	return tokens
}
