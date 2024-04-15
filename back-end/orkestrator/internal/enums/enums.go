package enums

import "github.com/Conty111/SuperCalculator/back-end/models"

type ApiType string

const (
	GrpcApi ApiType = "grpc"
	RestApi ApiType = "rest"
)

const (
	Common models.Role = "common"
	Admin  models.Role = "admin"
)
