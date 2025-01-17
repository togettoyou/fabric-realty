package api

import (
	"application/model"
	"application/service"
	"application/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: &service.UserService{},
	}
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.BadRequest(c, "用户注册信息格式错误")
		return
	}

	if err := h.userService.Register(&user); err != nil {
		utils.ServerError(c, "用户注册失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "用户注册成功", nil)
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		utils.BadRequest(c, "登录信息格式错误")
		return
	}

	user, token, err := h.userService.Login(loginData.Username, loginData.Password)
	if err != nil {
		utils.Unauthorized(c, "登录失败："+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"type":     user.Type,
			"name":     user.Name,
			"phone":    user.Phone,
			"email":    user.Email,
		},
	})
}
