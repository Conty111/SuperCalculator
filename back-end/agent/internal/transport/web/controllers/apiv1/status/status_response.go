package status

import (
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/app/build"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/models"
)

// Response is a declaration for a status response
type Response struct {
	ID     string        `jsonapi:"primary,status"`
	Status string        `jsonapi:"attr,status"`
	Info   *models.Stats `json:"info"`
	Build  *build.Info   `jsonapi:"attr,build"`
}
