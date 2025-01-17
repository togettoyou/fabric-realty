package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	MySQL struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"mysql"`

	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`

	JWT struct {
		Secret     string `yaml:"secret"`
		ExpireTime int    `yaml:"expire_time"`
	} `yaml:"jwt"`

	Fabric struct {
		ChannelName   string `yaml:"channelid"`
		ChaincodeName string `yaml:"chaincodeid"`
		Org1MSPID     string `yaml:"mspid"`
		CryptoPath    string `yaml:"cryptopath"`
		CertPath      string `yaml:"certpath"`
		KeyPath       string `yaml:"keypath"`
		TLSCertPath   string `yaml:"tlscertpath"`
		PeerEndpoint  string `yaml:"peerendpoint"`
		GatewayPeer   string `yaml:"gatewaypeer"`
	} `yaml:"fabric"`
}

var GlobalConfig Config

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.Unmarshal(&GlobalConfig)
}
