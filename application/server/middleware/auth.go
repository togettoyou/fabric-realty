package middleware

import (
	"application/model"
	"application/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "请求头中缺少Authorization认证信息")
			c.Abort()
			return
		}

		// 检查 Authorization header 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.Unauthorized(c, "认证信息格式错误，请使用Bearer {token}格式")
			c.Abort()
			return
		}

		// 解析 token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.Unauthorized(c, "无效或已过期的访问令牌")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_type", claims.UserType)

		c.Next()
	}
}

// RequireRoles 角色验证中间件
func RequireRoles(allowedTypes ...model.UserType) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			utils.Unauthorized(c, "未找到用户角色信息")
			c.Abort()
			return
		}

		allowed := false
		for _, t := range allowedTypes {
			if userType.(model.UserType) == t {
				allowed = true
				break
			}
		}

		if !allowed {
			utils.Fail(c, http.StatusForbidden, "您没有权限执行此操作")
			c.Abort()
			return
		}

		c.Next()
	}
}
