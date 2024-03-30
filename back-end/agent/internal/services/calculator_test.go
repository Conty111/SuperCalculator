package services_test

import (
	"github.com/Conty111/SuperCalculator/back-end/models"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Conty111/SuperCalculator/back-end/agent/internal/agent_errors"
	. "github.com/Conty111/SuperCalculator/back-end/agent/internal/services"
)

var _ = Describe("CalculatorService", func() {
	var es *CalculatorService

	BeforeEach(func() {
		es = NewCalculatorService()
	})

	Context("Calculate", func() {
		It("should calculate simple expression", func() {
			result, err := es.Calculate("2+3")
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(5.0))
		})
		It("should calculate expression with all operations", func() {
			result, err := es.Calculate("(2+3)*2/4 - 5")
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(float64(((2.0+3.0)*2.0)/4.0) - 5.0))
		})
		It("should calculate expression with floats", func() {
			result, err := es.Calculate("(2.12+3.0)*2")
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal((2.12 + 3) * 2))
		})
		It("should calculate expression with floats", func() {
			result, err := es.Calculate("(2.12+3.0)*2")
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal((2.12 + 3.0) * 2))
		})
		It("should calculate expression with whitespaces", func() {
			result, err := es.Calculate("2 + 2 * 3 / 3")
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(4.0))
		})
		It("should return error division by zero", func() {
			_, err := es.Calculate("2 / (1+ 1 -1-1)")
			Expect(err).To(HaveOccurred())
		})
		It("should return error mismatched parentheses", func() {
			_, err := es.Calculate("2 / (1+1-1")
			Expect(err).To(HaveOccurred())
		})
		It("should return error not enough operators or operand", func() {
			_, err := es.Calculate("2 / (1+ 1 -1)+")
			Expect(err).To(HaveOccurred())
		})
		It("should return error in subexpression", func() {
			_, err := es.Calculate("2+(1/0)")
			Expect(err).To(HaveOccurred())
		})
	})
	Context("Calculate with delays", func() {
		It("should calculate expression after delay with 0.3s delta", func() {
			es.SetOperationDuration(
				&models.DurationSettings{
					DivisionDuration: 200,
					SubtractDuration: 1000,
					MultiplyDuration: 1000,
					AddDuration:      1000,
				},
			)
			duration := es.AddTime.Milliseconds() + es.MultiplyTime.Milliseconds()
			ch := make(chan float64)
			t1 := time.Now()
			go func() {
				defer close(ch)
				result, err := es.Calculate("(2+3)*4")
				Expect(err).ToNot(HaveOccurred())
				ch <- result
			}()
			res := <-ch
			t2 := time.Since(t1)
			Expect(res).To(Equal(20.0))
			Expect(t2.Milliseconds()).To(BeNumerically(">=", duration-150))
			Expect(t2.Milliseconds()).To(BeNumerically("<=", duration+150))
		})
	})

	Context("Execute", func() {
		It("should execute the task and return result", func() {
			task := &models.Task{
				Expression: "1+2+3",
			}
			result := es.Execute(task)
			Expect(result).ToNot(BeNil())
			Expect(result.Error).To(Equal(""))
			Expect(result.Value).To(Equal(6.0))
		})
		It("should execute the task and return result with error", func() {
			task := &models.Task{
				Expression: "1+2+3/0",
			}
			result := es.Execute(task)
			Expect(result).ToNot(BeNil())
			Expect(result.Error).ToNot(Equal(""))
		})
		It("should execute the task and return result with error", func() {
			task := &models.Task{
				Expression: "1+2+3$0",
			}
			result := es.Execute(task)
			Expect(result).ToNot(BeNil())
			Expect(result.Error).ToNot(Equal(""))
		})
	})

	Context("ValidateExpression", func() {
		Context("Positive tests", func() {
			It("should validate valid expression", func() {
				validExpression := "(2+3)*4"
				validatedExpression, err := es.ValidateExpression(validExpression)
				Expect(err).ToNot(HaveOccurred())
				Expect(validatedExpression).To(Equal(validExpression))
			})
			It("should validate valid expression with whitespaces", func() {
				validExpression := " ( 2 +3) * 4"
				validatedExpression, err := es.ValidateExpression(validExpression)
				Expect(err).ToNot(HaveOccurred())
				Expect(validatedExpression).To(Equal(strings.ReplaceAll(validExpression, " ", "")))
			})
		})
		Context("with mismatched parentheses", func() {
			It("should return error", func() {
				invalidExpression := "(2+3)*4)"
				_, err := es.ValidateExpression(invalidExpression)
				Expect(err).To(HaveOccurred())
			})
		})
		Context("with invalid characters", func() {
			It("should return not allowed char", func() {
				invalidExpression := "(2+3)*$4"
				_, err := es.ValidateExpression(invalidExpression)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(agent_errors.NotAllowedChar))
			})
		})
		Context("with invalid float operand", func() {
			It("should return error", func() {
				invalidExpression := "2. + 2.2"
				_, err := es.ValidateExpression(invalidExpression)
				Expect(err).To(HaveOccurred())
			})
			It("should return error", func() {
				invalidExpression := "2.. + 2.2"
				_, err := es.ValidateExpression(invalidExpression)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
