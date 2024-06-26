package status_test

import (
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/app/build"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/app/initializers"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/Conty111/SuperCalculator/back-end/agent/internal/transport/web/controllers/apiv1/status"
)

var _ = Describe("Controller", func() {
	var (
		statusCtrl *Controller
		agentID    int32
		agentName  string
	)

	BeforeEach(func() {
		//gin.SetMode(gin.ReleaseMode)

		agentID = 0
		agentName = "Name"

		info := build.NewInfo()
		statusCtrl = NewController(info, initializers.InitializeMonitor(agentID, agentName))
	})

	It("controller should not be nil", func() {
		Expect(statusCtrl).NotTo(BeNil())
	})

	Describe("GetStatus()", func() {
		Context("default request", func() {
			It("should return status", func() {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)
				ctx.Request, _ = http.NewRequest("GET", "/api/v1/status", nil)

				statusCtrl.GetStatus(ctx)

				Expect(w.Result().StatusCode).To(Equal(http.StatusOK))
			})
		})
	})
})
