package health

import "github.com/gin-gonic/gin"

func (s *Service) HealthCheck(c *gin.Context) {

	c.JSON(200, gin.H{
		"healthy": true,
	})
}
