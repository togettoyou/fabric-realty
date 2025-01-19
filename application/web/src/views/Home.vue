<template>
  <div class="home">
    <div class="page-header">
      <h1 class="page-title">基于区块链的房地产交易系统</h1>
      <p class="page-description">选择您要进入的功能模块</p>
    </div>

    <div class="content">
      <router-link to="/realty-agency" class="card">
        <HomeOutlined class="card-icon" />
        <div class="card-title">房产管理</div>
        <div class="card-description">负责房产信息的登记和管理</div>
      </router-link>

      <router-link to="/trading-platform" class="card">
        <ShopOutlined class="card-icon" />
        <div class="card-title">交易管理</div>
        <div class="card-description">负责房产交易的创建和管理</div>
      </router-link>

      <router-link to="/bank" class="card">
        <CheckCircleOutlined class="card-icon" />
        <div class="card-title">交易确认</div>
        <div class="card-description">负责房产交易的确认和完成</div>
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { HomeOutlined, ShopOutlined, CheckCircleOutlined } from '@ant-design/icons-vue';
</script>

<style scoped>
.home {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #f0f2f5;
}

.content {
  flex: 1;
  display: flex;
  gap: 24px;
  padding: 24px;
}

.card {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 24px;
  border-radius: 8px;
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s;
  cursor: pointer;
  text-decoration: none;
  color: rgba(0, 0, 0, 0.85);
}

.card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
  transform: translateY(-2px);
}

.card-icon {
  font-size: 48px;
  margin-bottom: 16px;
  color: #1890ff;
}

.card-title {
  font-size: 18px;
  font-weight: 500;
  margin-bottom: 8px;
  color: rgba(0, 0, 0, 0.85);
}

.card-description {
  font-size: 14px;
  color: rgba(0, 0, 0, 0.45);
  text-align: center;
}

.page-header {
  padding: 24px;
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  text-align: center;
}

.page-title {
  font-size: 28px;
  font-weight: 500;
  color: rgba(0, 0, 0, 0.85);
  margin: 0;
}

.page-description {
  font-size: 16px;
  color: rgba(0, 0, 0, 0.45);
  margin: 8px 0 0 0;
}

@media (max-width: 768px) {
  .content {
    flex-direction: column;
  }
}
</style> 