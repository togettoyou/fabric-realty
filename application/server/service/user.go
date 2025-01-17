package service

import (
	"application/model"
	"application/utils"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

// Register 用户注册
func (s *UserService) Register(user *model.User) error {
	// 验证用户类型
	switch user.Type {
	case model.Buyer, model.Seller, model.RealtyAdmin, model.BankAdmin:
		// 合法的用户类型
	default:
		return fmt.Errorf("无效的用户类型：%s", user.Type)
	}

	// 管理员账号不允许注册
	if user.Type == model.RealtyAdmin || user.Type == model.BankAdmin {
		return fmt.Errorf("管理员账号不允许通过此接口注册")
	}

	// 检查用户名是否已存在
	var existingUser model.User
	if err := model.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return fmt.Errorf("用户名已存在")
	}

	// 验证必填字段
	if user.Username == "" || user.Password == "" || user.Name == "" || user.Phone == "" {
		return fmt.Errorf("用户名、密码、姓名和电话为必填项")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败：%v", err)
	}
	user.Password = string(hashedPassword)

	// 创建用户
	if err := model.DB.Create(user).Error; err != nil {
		return fmt.Errorf("创建用户失败：%v", err)
	}

	return nil
}

// Login 用户登录
func (s *UserService) Login(username, password string) (*model.User, string, error) {
	var user model.User
	if err := model.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, "", fmt.Errorf("用户不存在")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", fmt.Errorf("密码错误")
	}

	// 生成 JWT token
	token, err := utils.GenerateToken(&user)
	if err != nil {
		return nil, "", fmt.Errorf("生成令牌失败：%v", err)
	}

	return &user, token, nil
}
