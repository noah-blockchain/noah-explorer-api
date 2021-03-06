package validators

import "github.com/gin-gonic/gin"

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	validators := r.Group("/validators")
	{
		validators.GET("", GetAggregatedValidators)
		validators.GET("/:publicKey/transactions", GetValidatorTransactions)
		validators.GET("/:publicKey", GetValidator)
		validators.GET("/:publicKey/delegators", GetDelegators)
		//validators.GET("/ull", GetValidatorsFull)
	}
}
