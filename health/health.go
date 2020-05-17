package health

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckHealthHandler(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"health": "Rafiki is Healthy"})
}
