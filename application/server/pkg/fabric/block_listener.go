package fabric

import (
	"context"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-protos-go-apiv2/common"
)

// BlockEventListener 区块事件监听器
type BlockEventListener struct {
	sync.RWMutex
	networks map[string]*client.Network
	ctx      context.Context
	cancel   context.CancelFunc
	dataDir  string
}

var (
	listener     *BlockEventListener
	listenerOnce sync.Once
)

// InitBlockListener 初始化区块监听器
func InitBlockListener(dataDir string) {
	listenerOnce.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		listener = &BlockEventListener{
			networks: make(map[string]*client.Network),
			ctx:      ctx,
			cancel:   cancel,
			dataDir:  dataDir,
		}

		// 创建数据目录
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			fmt.Printf("创建数据目录失败：%v\n", err)
			return
		}
	})
}

// AddNetwork 添加网络
func AddNetwork(orgName string, network *client.Network) error {
	if listener == nil {
		return fmt.Errorf("区块监听器未初始化")
	}

	listener.Lock()
	defer listener.Unlock()

	// 为每个组织创建事件通道
	listener.networks[orgName] = network

	// 启动区块监听
	go listener.startBlockListener(orgName)

	return nil
}

// startBlockListener 启动区块监听
func (l *BlockEventListener) startBlockListener(orgName string) {
	network := l.networks[orgName]
	if network == nil {
		fmt.Printf("组织[%s]的网络未找到\n", orgName)
		return
	}

	// 获取当前已保存的最新区块号
	lastBlockNum := l.getLastBlockNum(orgName)
	startBlock := lastBlockNum + 1

	// 创建区块事件请求
	events, err := network.BlockEvents(l.ctx, client.WithStartBlock(startBlock))
	if err != nil {
		fmt.Printf("创建区块事件请求失败：%v\n", err)
		return
	}

	for {
		select {
		case <-l.ctx.Done():
			return
		case block := <-events:
			// 保存区块数据
			l.saveBlock(orgName, block)
		}
	}
}

// getLastBlockNum 获取最后保存的区块号
func (l *BlockEventListener) getLastBlockNum(orgName string) uint64 {
	filePath := filepath.Join(l.dataDir, fmt.Sprintf("%s_last_block.json", orgName))
	data, err := os.ReadFile(filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("读取最后区块号失败：%v\n", err)
		}
		return 0
	}

	var lastBlock struct {
		BlockNum uint64 `json:"block_num"`
	}
	if err := json.Unmarshal(data, &lastBlock); err != nil {
		fmt.Printf("解析最后区块号失败：%v\n", err)
		return 0
	}

	return lastBlock.BlockNum
}

// saveBlock 保存区块数据
func (l *BlockEventListener) saveBlock(orgName string, block *common.Block) {
	blockNum := block.GetHeader().GetNumber()

	// 计算区块哈希
	blockHeader := struct {
		Number       *big.Int
		PreviousHash []byte
		DataHash     []byte
	}{
		PreviousHash: block.GetHeader().GetPreviousHash(),
		DataHash:     block.GetHeader().GetDataHash(),
		Number:       new(big.Int).SetUint64(blockNum),
	}
	headerBytes, err := asn1.Marshal(blockHeader)
	if err != nil {
		fmt.Printf("序列化区块头失败：%v\n", err)
		return
	}
	blockHash := sha256.Sum256(headerBytes)

	// 保存区块数据
	blockData := struct {
		BlockNum  uint64 `json:"block_num"`
		BlockHash string `json:"block_hash"`
		DataHash  string `json:"data_hash"`
		PrevHash  string `json:"prev_hash"`
		TxCount   int    `json:"tx_count"`
	}{
		BlockNum:  blockNum,
		BlockHash: fmt.Sprintf("%x", blockHash[:]),
		DataHash:  fmt.Sprintf("%x", block.GetHeader().GetDataHash()),
		PrevHash:  fmt.Sprintf("%x", block.GetHeader().GetPreviousHash()),
		TxCount:   len(block.GetData().GetData()),
	}

	// 保存区块数据
	blockFile := filepath.Join(l.dataDir, fmt.Sprintf("%s_block_%d.json", orgName, blockNum))
	blockJSON, err := json.MarshalIndent(blockData, "", "  ")
	if err != nil {
		fmt.Printf("序列化区块数据失败：%v\n", err)
		return
	}
	if err := os.WriteFile(blockFile, blockJSON, 0644); err != nil {
		fmt.Printf("保存区块数据失败：%v\n", err)
		return
	}

	// 更新最后区块号
	lastBlockFile := filepath.Join(l.dataDir, fmt.Sprintf("%s_last_block.json", orgName))
	lastBlockJSON, err := json.Marshal(struct {
		BlockNum uint64    `json:"block_num"`
		SaveTime time.Time `json:"save_time"`
	}{
		BlockNum: blockNum,
		SaveTime: time.Now(),
	})
	if err != nil {
		fmt.Printf("序列化最后区块号失败：%v\n", err)
		return
	}
	if err := os.WriteFile(lastBlockFile, lastBlockJSON, 0644); err != nil {
		fmt.Printf("保存最后区块号失败：%v\n", err)
		return
	}

	fmt.Printf("已保存组织[%s]的区块[%d]\n", orgName, blockNum)
}
