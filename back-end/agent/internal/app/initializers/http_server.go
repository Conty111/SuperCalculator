package initializers

import (
	"fmt"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// InitializeHTTPServer create new http.Server instance
func InitializeHTTPServer(router *gin.Engine, cfg *config.HTTPConfig) *http.Server {
	// create http server
	return &http.Server{
		Addr:              fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler:           router,
		ReadHeaderTimeout: 30 * time.Second,
	}
}
