package status

import (
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/app/build"
)

// Response is a declaration for a status response
type Response struct {
	ID     string      `jsonapi:"primary,status"`
	Status string      `jsonapi:"attr,status"`
	Build  *build.Info `jsonapi:"attr,build"`
}
