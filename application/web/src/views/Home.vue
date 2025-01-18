<template>
  <div class="home">
    <a-row :gutter="[24, 24]" justify="center" align="middle" style="min-height: 100vh">
      <a-col :span="24">
        <h1 class="title">基于 Hyperledger Fabric 的房地产交易系统</h1>
      </a-col>
      <a-col :span="8">
        <a-card hoverable @click="$router.push('/realty-agency')">
          <template #cover>
            <img alt="不动产登记机构" src="../assets/realty.svg" />
          </template>
          <a-card-meta title="不动产登记机构">
            <template #description>负责房产信息的登记和所有权变更</template>
          </a-card-meta>
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card hoverable @click="$router.push('/trading-platform')">
          <template #cover>
            <img alt="交易平台" src="../assets/trading.svg" />
          </template>
          <a-card-meta title="交易平台">
            <template #description>负责创建和管理交易信息</template>
          </a-card-meta>
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card hoverable @click="$router.push('/bank')">
          <template #cover>
            <img alt="银行" src="../assets/bank.svg" />
          </template>
          <a-card-meta title="银行">
            <template #description>负责交易的完成确认</template>
          </a-card-meta>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import 'ant-design-vue/dist/antd.css';
</script>

<style scoped>
.home {
  padding: 24px;
  background-color: #f0f2f5;
}

.title {
  text-align: center;
  color: #1890ff;
  font-size: 2.5em;
  margin-bottom: 48px;
}

:deep(.ant-card) {
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
}

:deep(.ant-card:hover) {
  transform: translateY(-6px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
}

:deep(.ant-card-cover img) {
  padding: 24px;
  height: 200px;
  object-fit: contain;
}
</style> 