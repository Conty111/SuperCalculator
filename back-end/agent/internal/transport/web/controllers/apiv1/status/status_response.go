package status

import (
	"github.com/Conty111/SuperCalculator/back-end/models"
)

// Response is a declaration for a status response
type Response struct {
	ID   int32             `json:"id"`
	Info *models.AgentInfo `json:"info"`
}
