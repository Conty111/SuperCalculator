package initializers

import "github.com/Conty111/SuperCalculator/back-end/agent/internal/services"

func InitializeExpressionService() *services.CalculatorService {
	return services.NewCalculatorService()
}
