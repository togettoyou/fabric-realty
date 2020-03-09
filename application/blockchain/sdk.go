/**
 * @Author: 夜央 Oh oh oh oh oh oh (https://github.com/togettoyou)
 * @Email: zoujh99@qq.com
 * @Date: 2020/3/4 4:21 下午
 * @Description: fabric-sdk-go的封装
 */
package blockchain

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// 配置信息
var (
	SDK           *fabsdk.FabricSDK          // Fabric提供的SDK
	ChannelName   = "assetschannel"          // 通道名称
	ChainCodeName = "blockchain-real-estate" // 链码名称
	Org           = "org1"                   // 组织名称
	User          = "Admin"                  // 用户
	ConfigPath    = "blockchain/config.yaml" // 配置文件路径
)

// 初始化
func Init() {
	var err error
	// 通过配置文件初始化SDK
	SDK, err = fabsdk.New(config.FromFile(ConfigPath))
	if err != nil {
		panic(err)
	}
}

// 区块链交互
func ChannelExecute(fcn string, args [][]byte) (channel.Response, error) {
	// 创建客户端，表明在通道的身份
	ctx := SDK.ChannelContext(ChannelName, fabsdk.WithOrg(Org), fabsdk.WithUser(User))
	cli, err := channel.New(ctx)
	if err != nil {
		return channel.Response{}, err
	}
	// 对区块链增删改的操作（调用了链码的invoke）
	resp, err := cli.Execute(channel.Request{
		ChaincodeID: ChainCodeName,
		Fcn:         fcn,
		Args:        args,
	}, channel.WithTargetEndpoints("peer0.org1.blockchainrealestate.com"))
	if err != nil {
		return channel.Response{}, err
	}
	//返回链码执行后的结果
	return resp, nil
}

// 区块链查询
func ChannelQuery(fcn string, args [][]byte) (channel.Response, error) {
	// 创建客户端，表明在通道的身份
	ctx := SDK.ChannelContext(ChannelName, fabsdk.WithOrg(Org), fabsdk.WithUser(User))
	cli, err := channel.New(ctx)
	if err != nil {
		return channel.Response{}, err
	}
	// 对区块链查询的操作（调用了链码的invoke），将结果返回
	resp, err := cli.Query(channel.Request{
		ChaincodeID: ChainCodeName,
		Fcn:         fcn,
		Args:        args,
	}, channel.WithTargetEndpoints("peer0.org1.blockchainrealestate.com"))
	if err != nil {
		return channel.Response{}, err
	}
	//返回链码执行后的结果
	return resp, nil
}
