package model

import (
	"fmt"
	"log"

	"application/config"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.GlobalConfig.MySQL.User,
		config.GlobalConfig.MySQL.Password,
		config.GlobalConfig.MySQL.Host,
		config.GlobalConfig.MySQL.Port,
		config.GlobalConfig.MySQL.DBName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接数据库失败：%v", err)
	}

	// 自动迁移数据库结构
	err = DB.AutoMigrate(&User{})
	if err != nil {
		return fmt.Errorf("数据库迁移失败：%v", err)
	}

	// 创建默认管理员账户
	createDefaultAdmins()

	return nil
}

func createDefaultAdmins() {
	password, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)

	// 创建默认房管局管理员
	realtyAdmin := User{
		Username: "realty_admin",
		Password: string(password),
		Type:     RealtyAdmin,
		Name:     "房管局管理员",
	}

	// 创建默认银行管理员
	bankAdmin := User{
		Username: "bank_admin",
		Password: string(password),
		Type:     BankAdmin,
		Name:     "银行管理员",
	}

	// 如果不存在则创建
	if err := DB.FirstOrCreate(&realtyAdmin, User{Username: "realty_admin"}).Error; err != nil {
		log.Printf("创建默认房管局管理员失败：%v", err)
	}

	if err := DB.FirstOrCreate(&bankAdmin, User{Username: "bank_admin"}).Error; err != nil {
		log.Printf("创建默认银行管理员失败：%v", err)
	}
}
