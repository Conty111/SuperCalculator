package initializers

import "github.com/Conty111/SuperCalculator/back-end/agent/internal/services"

func InitializeExpressionService() *services.ExpressionService {
	return services.NewExpressionService()
}
