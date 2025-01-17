package model

import (
	"gorm.io/gorm"
)

type UserType string

const (
	Buyer       UserType = "buyer"        // 买家
	Seller      UserType = "seller"       // 卖家
	RealtyAdmin UserType = "realty_admin" // 房管局管理员
	BankAdmin   UserType = "bank_admin"   // 银行管理员
)

// User 用户模型
type User struct {
	gorm.Model
	Username string   `gorm:"type:varchar(255);uniqueIndex;not null" json:"username"` // 用户名
	Password string   `gorm:"type:varchar(255);not null" json:"password"`             // 密码
	Type     UserType `gorm:"type:varchar(20);not null" json:"type"`                  // 用户类型
	Name     string   `gorm:"type:varchar(50)" json:"name"`                           // 真实姓名
	Phone    string   `gorm:"type:varchar(20)" json:"phone"`                          // 电话
	Email    string   `gorm:"type:varchar(100)" json:"email"`                         // 邮箱
}
