package initializers_test

import (
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/app/dependencies"
	. "github.com/Conty111/SuperCalculator/back-end/agent/internal/app/initializers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Router", func() {
	Describe("InitializeRouter()", func() {
		var (
			c *dependencies.Container
		)

		BeforeEach(func() {
			c = &dependencies.Container{}
		})

		It("should initialize router", func() {
			r := InitializeRouter(c)

			Expect(r).NotTo(BeNil())
		})
	})
})
