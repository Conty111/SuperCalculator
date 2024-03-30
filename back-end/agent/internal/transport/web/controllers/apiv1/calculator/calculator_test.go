package calculator

//
//import (
//	. "github.com/onsi/ginkgo/v2"
//	. "github.com/onsi/gomega"
//
//	. "github.com/Conty111/SuperCalculator/back-end/agent/internal/services"
//)
//
//var _ = Describe("CalculatorService", func() {
//	var es *CalculatorService
//
//	BeforeEach(func() {
//		es = NewCalculatorService()
//	})
//
//	Context("Calculate", func() {
//		It("should calculate simple expression", func() {
//			result, err := es.Calculate("2+3")
//			Expect(err).ToNot(HaveOccurred())
//			Expect(result).To(Equal(5.0))
//		})
//		It("should calculate expression with all operations", func() {
//			result, err := es.Calculate("(2+3)*2/4 - 5")
//			Expect(err).ToNot(HaveOccurred())
//			Expect(result).To(Equal(float64(((2.0+3.0)*2.0)/4.0) - 5.0))
//		})
//	})
//})
