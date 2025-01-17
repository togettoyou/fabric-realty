package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 配置
type Config struct {
	MySQL  MySQLConfig  `yaml:"mysql"`
	Server ServerConfig `yaml:"server"`
	JWT    JWTConfig    `yaml:"jwt"`
	Fabric FabricConfig `yaml:"fabric"`
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int `yaml:"port"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string `yaml:"secret"`
	ExpireTime int64  `yaml:"expire_time"`
}

// FabricConfig Fabric配置
type FabricConfig struct {
	ChannelName   string                        `yaml:"channelName"`
	ChaincodeName string                        `yaml:"chaincodeName"`
	Organizations map[string]OrganizationConfig `yaml:"organizations"`
}

// OrganizationConfig 组织配置
type OrganizationConfig struct {
	MSPID        string `yaml:"mspID"`
	CertPath     string `yaml:"certPath"`
	KeyPath      string `yaml:"keyPath"`
	TLSCertPath  string `yaml:"tlsCertPath"`
	PeerEndpoint string `yaml:"peerEndpoint"`
	GatewayPeer  string `yaml:"gatewayPeer"`
}

var GlobalConfig Config

// InitConfig 初始化配置
func InitConfig() error {
	// 读取配置文件
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return fmt.Errorf("读取配置文件失败：%v", err)
	}

	// 解析配置文件
	err = yaml.Unmarshal(data, &GlobalConfig)
	if err != nil {
		return fmt.Errorf("解析配置文件失败：%v", err)
	}

	return nil
}
