package model

import (
	"fmt"
	"log"

	"application/config"

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
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// 自动迁移数据库结构
	err = DB.AutoMigrate(&User{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	// 创建默认管理员账户
	createDefaultAdmins()

	return nil
}

func createDefaultAdmins() {
	// 创建默认房管局管理员
	realtyAdmin := User{
		Username: "realty_admin",
		Password: "123456", // 实际应用中需要加密
		Type:     RealtyAdmin,
		Name:     "房管局管理员",
	}

	// 创建默认银行管理员
	bankAdmin := User{
		Username: "bank_admin",
		Password: "123456", // 实际应用中需要加密
		Type:     BankAdmin,
		Name:     "银行管理员",
	}

	// 如果不存在则创建
	if err := DB.FirstOrCreate(&realtyAdmin, User{Username: "realty_admin"}).Error; err != nil {
		log.Printf("Failed to create default realty admin: %v", err)
	}

	if err := DB.FirstOrCreate(&bankAdmin, User{Username: "bank_admin"}).Error; err != nil {
		log.Printf("Failed to create default bank admin: %v", err)
	}
}
