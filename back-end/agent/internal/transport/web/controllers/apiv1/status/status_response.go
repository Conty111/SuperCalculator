package status

import (
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/models"
)

// Response is a declaration for a status response
type Response struct {
	ID     int32         `json:"id"`
	Status string        `json:"status"`
	Info   *models.Stats `json:"info"`
	//Build  *build.Info   `json:"build"`
}
