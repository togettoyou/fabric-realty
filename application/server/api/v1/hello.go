package v1

import (
	"application/pkg/app"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Hello(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, "成功", map[string]interface{}{
		"msg": "Hello",
	})
}
