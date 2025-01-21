package fabric

import (
	"application/config"
	"crypto/x509"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

var (
	// 组织对应的合约客户端
	contracts = make(map[string]*client.Contract)
)

// InitFabric 初始化 Fabric 客户端
func InitFabric() error {
	// 初始化区块监听器
	if err := initBlockListener(filepath.Join("data", "blocks")); err != nil {
		return fmt.Errorf("初始化区块监听器失败: %w", err)
	}

	// 为每个组织创建合约客户端
	for orgName, orgConfig := range config.GlobalConfig.Fabric.Organizations {
		// 创建 gRPC 连接
		clientConnection, err := newGrpcConnection(orgConfig)
		if err != nil {
			return fmt.Errorf("创建组织[%s]的gRPC连接失败：%v", orgName, err)
		}

		// 创建组织身份
		id, err := newIdentity(orgConfig)
		if err != nil {
			return fmt.Errorf("创建组织[%s]身份失败：%v", orgName, err)
		}

		// 创建签名函数
		sign, err := newSign(orgConfig)
		if err != nil {
			return fmt.Errorf("创建组织[%s]签名函数失败：%v", orgName, err)
		}

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
			return fmt.Errorf("连接组织[%s]的Fabric网关失败：%v", orgName, err)
		}

		network := gw.GetNetwork(config.GlobalConfig.Fabric.ChannelName)
		contracts[orgName] = network.GetContract(config.GlobalConfig.Fabric.ChaincodeName)

		// 添加网络到区块监听器
		if err := addNetwork(orgName, network); err != nil {
			return fmt.Errorf("添加网络到区块监听器失败：%v", err)
		}
	}

	return nil
}

// GetContract 获取指定组织的合约客户端
func GetContract(orgName string) *client.Contract {
	return contracts[orgName]
}

// ExtractErrorMessage 从错误中提取详细信息
func ExtractErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	// 尝试获取 gRPC 状态
	if st, ok := status.FromError(err); ok {
		// 获取详细信息
		msg := st.Message()
		details := st.Details()
		code := st.Code()

		// 构建完整的错误信息
		fullError := fmt.Sprintf("错误码: %v, 消息: %v", code, msg)
		if len(details) > 0 {
			fullError += fmt.Sprintf(", 详情: %+v", details)
		}
		return fullError
	}
	return err.Error()
}

// newGrpcConnection 创建 gRPC 连接
func newGrpcConnection(orgConfig config.OrganizationConfig) (*grpc.ClientConn, error) {
	certificatePEM, err := os.ReadFile(orgConfig.TLSCertPath)
	if err != nil {
		return nil, fmt.Errorf("读取TLS证书文件失败：%w", err)
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		return nil, fmt.Errorf("解析TLS证书失败：%w", err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, orgConfig.GatewayPeer)

	connection, err := grpc.Dial(orgConfig.PeerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		return nil, fmt.Errorf("创建gRPC连接失败：%w", err)
	}

	return connection, nil
}

// newIdentity 创建身份
func newIdentity(orgConfig config.OrganizationConfig) (*identity.X509Identity, error) {
	certificatePEM, err := readFirstFile(orgConfig.CertPath)
	if err != nil {
		return nil, fmt.Errorf("读取证书文件失败：%w", err)
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		return nil, err
	}

	id, err := identity.NewX509Identity(orgConfig.MSPID, certificate)
	if err != nil {
		return nil, err
	}

	return id, nil
}

// newSign 创建签名函数
func newSign(orgConfig config.OrganizationConfig) (identity.Sign, error) {
	privateKeyPEM, err := readFirstFile(orgConfig.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("读取私钥文件失败：%w", err)
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, err
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		return nil, err
	}

	return sign, nil
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
