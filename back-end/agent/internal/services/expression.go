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
	operands, operators, err := es.ParseToInfix(expression)
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

// ParseToInfix returns operands and operators from expression in preorder
func (es *ExpressionService) ParseToInfix(expression string) ([]float64, []string, error) {
	// Разбираем выражение на операторы и операнды
	operators := make([]string, 0)
	operands := make([]float64, 0)
	operatorStack := make([]string, 0)

	// Функция для обработки операторов в стеке с учетом приоритета
	processOperators := func(i int) {
		for len(operatorStack) > 0 &&
			es.getOperatorPriority(operatorStack[len(operatorStack)-1]) >= es.getOperatorPriority(string(expression[i])) {
			operators = append(operators, operatorStack[len(operatorStack)-1])
			operatorStack = operatorStack[:len(operatorStack)-1]
		}
	}

	for i := 0; i < len(expression); {
		switch char := expression[i]; char {
		case '+', '-', '*', '/', '^':
			processOperators(i)
			operatorStack = append(operatorStack, string(char))
			i++
		case '(', ')':
			if char == '(' {
				operatorStack = append(operatorStack, string(char))
			} else {
				for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1] != "(" {
					operators = append(operators, operatorStack[len(operatorStack)-1])
					operatorStack = operatorStack[:len(operatorStack)-1]
				}
				if len(operatorStack) == 0 {
					return nil, nil, fmt.Errorf("несогласованные скобки в выражении")
				}
				operatorStack = operatorStack[:len(operatorStack)-1] // Убираем открывающую скобку
			}
			i++
		default:
			// Читаем операнд
			var operandString string
			for ; i < len(expression) && (expression[i] == '.' || unicode.IsDigit(rune(expression[i]))); i++ {
				operandString += string(expression[i])
			}

			operand, err := strconv.ParseFloat(operandString, 64)
			if err != nil {
				return nil, nil, agent_errors.ErrOperandParsing(operandString, err)
			}
			operands = append(operands, operand)
		}
	}

	processOperators(len(expression) - 1)

	return operands, operators, nil
}

// getOperatorPriority возвращает приоритет оператора
func (es *ExpressionService) getOperatorPriority(operator string) int {
	switch operator {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	case "^":
		return 3
	default:
		return 0
	}
}

// CalculateInfix calculate preorder operands and operators
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
			time.Sleep(es.MultiplyTime * time.Duration(operands[i+1]))
			result = math.Pow(result, operands[i+1])
		}
	}
	return result
}
