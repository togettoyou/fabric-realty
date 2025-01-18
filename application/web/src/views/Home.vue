<template>
  <div class="home">
    <div class="main-content">
      <div class="hero-section">
        <h1 class="hero-title">基于区块链的房地产交易系统</h1>
        <p class="hero-subtitle">使用 Hyperledger Fabric 构建的安全、透明、高效的房地产交易平台</p>
        <div class="hero-features">
          <div class="feature-item">
            <CheckCircleOutlined />
            <span>安全可靠</span>
          </div>
          <div class="feature-item">
            <SafetyCertificateOutlined />
            <span>数据透明</span>
          </div>
          <div class="feature-item">
            <ThunderboltOutlined />
            <span>高效便捷</span>
          </div>
        </div>
      </div>

      <div class="roles-section">
        <a-row :gutter="[32, 32]" justify="center">
          <a-col :xs="24" :sm="24" :md="8">
            <a-card hoverable class="role-card" @click="$router.push('/realty-agency')">
              <div class="card-content">
                <div class="card-icon">
                  <img alt="不动产登记机构" src="../assets/realty.svg" />
                </div>
                <div class="card-info">
                  <h3>不动产登记机构</h3>
                  <p>负责房产信息的登记和所有权变更</p>
                </div>
              </div>
            </a-card>
          </a-col>
          <a-col :xs="24" :sm="24" :md="8">
            <a-card hoverable class="role-card" @click="$router.push('/trading-platform')">
              <div class="card-content">
                <div class="card-icon">
                  <img alt="交易平台" src="../assets/trading.svg" />
                </div>
                <div class="card-info">
                  <h3>交易平台</h3>
                  <p>负责创建和管理交易信息</p>
                </div>
              </div>
            </a-card>
          </a-col>
          <a-col :xs="24" :sm="24" :md="8">
            <a-card hoverable class="role-card" @click="$router.push('/bank')">
              <div class="card-content">
                <div class="card-icon">
                  <img alt="银行" src="../assets/bank.svg" />
                </div>
                <div class="card-info">
                  <h3>银行</h3>
                  <p>负责交易的完成确认</p>
                </div>
              </div>
            </a-card>
          </a-col>
        </a-row>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { CheckCircleOutlined, SafetyCertificateOutlined, ThunderboltOutlined } from '@ant-design/icons-vue';
</script>

<style scoped>
.home {
  height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.main-content {
  width: 100%;
  max-width: 1200px;
}

.hero-section {
  text-align: center;
  margin-bottom: 48px;
}

.hero-title {
  font-size: 2.5em;
  font-weight: bold;
  margin-bottom: 16px;
  background: linear-gradient(45deg, #1890ff 0%, #36cfc9 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1);
}

.hero-subtitle {
  font-size: 1.2em;
  color: #666;
  margin-bottom: 32px;
}

.hero-features {
  display: flex;
  justify-content: center;
  gap: 48px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #1890ff;
}

.feature-item :deep(.anticon) {
  font-size: 20px;
}

.role-card {
  height: 100%;
  transition: all 0.3s ease;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  background: white;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.role-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.1);
}

.card-content {
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 24px;
}

.card-icon {
  width: 64px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(45deg, #f5f7fa 0%, #e4e7eb 100%);
  border-radius: 12px;
  padding: 12px;
}

.card-icon img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  transition: all 0.3s ease;
}

.role-card:hover .card-icon img {
  transform: scale(1.1);
}

.card-info {
  flex: 1;
}

.card-info h3 {
  font-size: 1.2em;
  color: #1890ff;
  margin-bottom: 8px;
}

.card-info p {
  color: #666;
  margin: 0;
  font-size: 14px;
}

@media (max-width: 768px) {
  .home {
    height: auto;
    padding: 48px 24px;
  }
  
  .hero-title {
    font-size: 2em;
  }
  
  .hero-subtitle {
    font-size: 1em;
  }
  
  .hero-features {
    flex-direction: column;
    gap: 16px;
  }
  
  .card-content {
    flex-direction: column;
    text-align: center;
    gap: 16px;
  }
  
  .card-icon {
    margin: 0 auto;
  }
}
</style> 