<template>
  <div class="home">
    <div class="hero-section">
      <div class="hero-content">
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
    </div>

    <div class="roles-section">
      <a-row :gutter="[32, 32]" justify="center">
        <a-col :xs="24" :sm="24" :md="8">
          <a-card hoverable class="role-card" @click="$router.push('/realty-agency')">
            <template #cover>
              <div class="card-cover">
                <img alt="不动产登记机构" src="../assets/realty.svg" />
              </div>
            </template>
            <a-card-meta title="不动产登记机构">
              <template #description>
                <div class="card-description">
                  <p>负责房产信息的登记和所有权变更</p>
                  <ul>
                    <li>创建和管理房产信息</li>
                    <li>记录房产所有权变更</li>
                    <li>维护房产状态</li>
                  </ul>
                </div>
              </template>
            </a-card-meta>
          </a-card>
        </a-col>
        <a-col :xs="24" :sm="24" :md="8">
          <a-card hoverable class="role-card" @click="$router.push('/trading-platform')">
            <template #cover>
              <div class="card-cover">
                <img alt="交易平台" src="../assets/trading.svg" />
              </div>
            </template>
            <a-card-meta title="交易平台">
              <template #description>
                <div class="card-description">
                  <p>负责创建和管理交易信息</p>
                  <ul>
                    <li>发起房产交易</li>
                    <li>记录交易信息</li>
                    <li>跟踪交易状态</li>
                  </ul>
                </div>
              </template>
            </a-card-meta>
          </a-card>
        </a-col>
        <a-col :xs="24" :sm="24" :md="8">
          <a-card hoverable class="role-card" @click="$router.push('/bank')">
            <template #cover>
              <div class="card-cover">
                <img alt="银行" src="../assets/bank.svg" />
              </div>
            </template>
            <a-card-meta title="银行">
              <template #description>
                <div class="card-description">
                  <p>负责交易的完成确认</p>
                  <ul>
                    <li>确认交易付款</li>
                    <li>完成交易流程</li>
                    <li>更新交易状态</li>
                  </ul>
                </div>
              </template>
            </a-card-meta>
          </a-card>
        </a-col>
      </a-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { CheckCircleOutlined, SafetyCertificateOutlined, ThunderboltOutlined } from '@ant-design/icons-vue';
</script>

<style scoped>
.home {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.hero-section {
  padding: 80px 24px;
  background: linear-gradient(45deg, #1890ff 0%, #36cfc9 100%);
  text-align: center;
  color: white;
}

.hero-content {
  max-width: 800px;
  margin: 0 auto;
}

.hero-title {
  font-size: 3em;
  font-weight: bold;
  margin-bottom: 24px;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.2);
  background: linear-gradient(45deg, #ffffff 0%, #f0f2f5 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.hero-subtitle {
  font-size: 1.4em;
  margin-bottom: 48px;
  opacity: 0.9;
}

.hero-features {
  display: flex;
  justify-content: center;
  gap: 48px;
  margin-top: 48px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 1.2em;
}

.feature-item :deep(.anticon) {
  font-size: 24px;
}

.roles-section {
  padding: 64px 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.role-card {
  height: 100%;
  transition: all 0.3s ease;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.role-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.1);
}

.card-cover {
  background: linear-gradient(45deg, #f5f7fa 0%, #e4e7eb 100%);
  padding: 48px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.card-cover img {
  height: 120px;
  width: 120px;
  object-fit: contain;
  transition: all 0.3s ease;
}

.role-card:hover .card-cover img {
  transform: scale(1.1);
}

:deep(.ant-card-meta-title) {
  font-size: 1.4em;
  color: #1890ff;
  margin-bottom: 16px !important;
}

.card-description {
  color: #666;
}

.card-description p {
  margin-bottom: 12px;
  font-weight: 500;
}

.card-description ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.card-description li {
  position: relative;
  padding-left: 20px;
  margin-bottom: 8px;
  font-size: 14px;
}

.card-description li::before {
  content: "•";
  color: #1890ff;
  position: absolute;
  left: 0;
  font-weight: bold;
}

@media (max-width: 768px) {
  .hero-title {
    font-size: 2em;
  }
  
  .hero-subtitle {
    font-size: 1.2em;
  }
  
  .hero-features {
    flex-direction: column;
    gap: 24px;
  }
}
</style> 