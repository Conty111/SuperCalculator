package status_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/Conty111/SuperCalculator/back-end/agent/internal/app/build"
	"github.com/gin-gonic/gin"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/Conty111/SuperCalculator/back-end/agent/internal/transport/web/controllers/apiv1/status"
)

var _ = Describe("Controller", func() {
	var (
		statusCtrl *Controller
	)

	BeforeEach(func() {
		gin.SetMode(gin.ReleaseMode)

		info := build.NewInfo()
		statusCtrl = NewController(info)
	})

	It("controller should not be nil", func() {
		Expect(statusCtrl).NotTo(BeNil())
	})

	Describe("GetStatus()", func() {
		It("should return status", func() {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request, _ = http.NewRequest("GET", "/api/v1/status", nil)

			statusCtrl.GetStatus(ctx)

			Expect(w.Result().StatusCode).To(Equal(http.StatusOK))
		})
	})
})
