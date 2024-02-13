package services_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/Conty111/SuperCalculator/back-end/agent/internal/agent_errors"
	. "github.com/Conty111/SuperCalculator/back-end/agent/internal/services"
)

type ExpressionServiceSuite struct {
	suite.Suite
	es *ExpressionService
}

func (suite *ExpressionServiceSuite) SetupTest() {
	suite.es = NewExpressionService()
}

func (suite *ExpressionServiceSuite) TestCalculate() {
	suite.T().Run("Simple expression without delays", func(t *testing.T) {
		result, err := suite.es.Calculate("2+3")
		assert.NoError(t, err)
		assert.Equal(t, 5.0, result)
	})

	suite.T().Run("Expression with delays", func(t *testing.T) {
		suite.es.AddTime = 100 * time.Millisecond
		suite.es.MultiplyTime = 50 * time.Millisecond

		result, err := suite.es.Calculate("(2+3)*4")
		assert.NoError(t, err)
		assert.Equal(t, 20.0, result)
	})
}

func (suite *ExpressionServiceSuite) TestValidateExpression() {
	suite.T().Run("Valid expression", func(t *testing.T) {
		validExpression := "(2+3)*4"
		validatedExpression, err := suite.es.ValidateExpression(validExpression)
		assert.NoError(t, err)
		assert.Equal(t, validExpression, validatedExpression)
	})

	suite.T().Run("Expression with invalid characters", func(t *testing.T) {
		invalidExpression := "(2+3)*$4"
		_, err := suite.es.ValidateExpression(invalidExpression)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "недопустимый символ")
	})

	suite.T().Run("Mismatched parentheses", func(t *testing.T) {
		invalidExpression := "(2+3)*4)"
		_, err := suite.es.ValidateExpression(invalidExpression)
		assert.Error(t, err)
		assert.Equal(t, agent_errors.ErrParenthesisNotValid, err)
	})

	suite.T().Run("Valid expression with spaces", func(t *testing.T) {
		validExpressionWithSpaces := "( 2 + 3 ) * 4"
		validatedExpression, err := suite.es.ValidateExpression(validExpressionWithSpaces)
		assert.NoError(t, err)
		assert.Equal(t, "(2+3)*4", validatedExpression)
	})
}

func (suite *ExpressionServiceSuite) TestParseToInfix() {
	suite.T().Run("Simple expression", func(t *testing.T) {
		operands, operators, err := suite.es.ParseToInfix("2+3")
		assert.NoError(t, err)
		assert.Equal(t, []float64{2, 3}, operands)
		assert.Equal(t, []string{"+"}, operators)
	})

	suite.T().Run("Expression with parentheses", func(t *testing.T) {
		operands, operators, err := suite.es.ParseToInfix("(2+3)*4")
		assert.NoError(t, err)
		assert.Equal(t, []float64{5, 4}, operands)
		assert.Equal(t, []string{"*"}, operators)
	})

	suite.T().Run("Expression with multiple operators", func(t *testing.T) {
		operands, operators, err := suite.es.ParseToInfix("2+3*4")
		assert.NoError(t, err)
		assert.Equal(t, []float64{2, 3, 4}, operands)
		assert.Equal(t, []string{"+", "*"}, operators)
	})
}

func (suite *ExpressionServiceSuite) TestCalculateInfix() {
	suite.T().Run("Simple expression without delays", func(t *testing.T) {
		result := suite.es.CalculateInfix([]float64{2, 3}, []string{"+"})
		assert.Equal(t, 5.0, result)
	})

	suite.T().Run("Expression with multiplication delay", func(t *testing.T) {
		suite.es.MultiplyTime = 100 * time.Millisecond
		result := suite.es.CalculateInfix([]float64{2, 3, 4}, []string{"+", "*"})
		assert.Equal(t, 20.0, result)
	})
}

func TestExpressionServiceSuite(t *testing.T) {
	suite.Run(t, new(ExpressionServiceSuite))
}
