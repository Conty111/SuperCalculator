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
		suite.es.AddTime = 3 * time.Second
		suite.es.MultiplyTime = 50 * time.Millisecond

		result, err := suite.es.Calculate("2+3*4")
		assert.NoError(t, err)
		assert.Equal(t, 14.0, result)
	})
	suite.T().Run("Expression with float", func(t *testing.T) {
		result, err := suite.es.Calculate("(2.12+3.0)*2")
		assert.NoError(t, err)
		assert.Equal(t, 10.24, result)
	})
	suite.T().Run("Expression a?", func(t *testing.T) {
		res, err := suite.es.Calculate("-2 + 3 * 4 * 1")
		assert.NoError(t, err)
		assert.Equal(t, 10.0, res)
	})
	suite.T().Run("Delays", func(t *testing.T) {
		suite.es.AddTime = 1 * time.Second
		suite.es.MultiplyTime = 2 * time.Second
		ch := make(chan float64)
		t1 := time.Now()
		go func() {
			defer close(ch)
			result, err := suite.es.Calculate("(2+3)*4")
			assert.NoError(t, err)
			ch <- result
		}()
		res := <-ch
		t2 := time.Since(t1)
		assert.Equal(t, 20.0, res)
		assert.LessOrEqual(t, t2, time.Second*4)
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

	suite.T().Run("Expression with invalid float operand", func(t *testing.T) {
		validExpressionWithSpaces := "2. * 2"
		_, err := suite.es.ValidateExpression(validExpressionWithSpaces)
		assert.Error(t, err)
	})
}

func TestExpressionServiceSuite(t *testing.T) {
	suite.Run(t, new(ExpressionServiceSuite))
}
