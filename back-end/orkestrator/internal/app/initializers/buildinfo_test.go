package initializers_test

import (
	. "github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/app/initializers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Buildinfo", func() {
	Describe("InitializeBuildInfo", func() {
		It("should initialize and return build.Info", func() {
			info := InitializeBuildInfo()

			Expect(info).NotTo(BeNil())
		})
	})
})
