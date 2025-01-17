package utils

import (
	"application/config"
	"crypto/x509"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	Contract *client.Contract
)

// InitFabric 初始化 Fabric 客户端
func InitFabric() error {
	clientConnection := newGrpcConnection()
	id := newIdentity()
	sign := newSign()

	// 创建 Gateway 连接
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return fmt.Errorf("连接Fabric网关失败：%v", err)
	}

	network := gw.GetNetwork(config.GlobalConfig.Fabric.ChannelName)
	Contract = network.GetContract(config.GlobalConfig.Fabric.ChaincodeName)

	return nil
}

// newGrpcConnection 创建 gRPC 连接
func newGrpcConnection() *grpc.ClientConn {
	certificatePEM, err := os.ReadFile(config.GlobalConfig.Fabric.TLSCertPath)
	if err != nil {
		panic(fmt.Errorf("读取TLS证书文件失败：%w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, config.GlobalConfig.Fabric.GatewayPeer)

	connection, err := grpc.Dial(config.GlobalConfig.Fabric.PeerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("创建gRPC连接失败：%w", err))
	}

	return connection
}

// newIdentity 创建身份
func newIdentity() *identity.X509Identity {
	certificatePEM, err := readFirstFile(config.GlobalConfig.Fabric.CertPath)
	if err != nil {
		panic(fmt.Errorf("读取证书文件失败：%w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(config.GlobalConfig.Fabric.Org1MSPID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

// newSign 创建签名函数
func newSign() identity.Sign {
	privateKeyPEM, err := readFirstFile(config.GlobalConfig.Fabric.KeyPath)
	if err != nil {
		panic(fmt.Errorf("读取私钥文件失败：%w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

// readFirstFile 读取目录中的第一个文件
func readFirstFile(dirPath string) ([]byte, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}

	fileNames, err := dir.Readdirnames(1)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(path.Join(dirPath, fileNames[0]))
}
