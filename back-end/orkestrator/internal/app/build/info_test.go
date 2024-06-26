package build_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/app/build"
)

var _ = Describe("Info", func() {
	Describe("NewInfo()", func() {
		It("should create new info object", func() {
			info := NewInfo()

			Expect(info).NotTo(BeNil())
		})
	})
})
