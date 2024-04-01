package initializers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"

	. "github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/app/initializers"
	"github.com/gobuffalo/envy"
)

var _ = Describe("Envs", func() {
	Describe("InitializeEnvs()", func() {
		var (
			k, v string
		)

		BeforeEach(func() {
			k = "SOME_TEST_ENV"
			v = "SOME_TEST_ENV_VALUE"

			Expect(os.Setenv(k, v)).To(BeNil())
		})

		It("should initialize envs with Envy package", func() {
			InitializeEnvs()

			Expect(os.Getenv(k)).To(Equal(v))
			Expect(envy.Get(k, "")).To(Equal(v))
		})
	})
})
