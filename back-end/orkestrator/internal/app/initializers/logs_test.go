package initializers_test

import (
	. "github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/app/initializers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logs", func() {
	Describe("InitializeLogs() ", func() {
		It("should initialize logs", func() {
			err := InitializeLogs()

			Expect(err).To(BeNil())
		})
	})
})
