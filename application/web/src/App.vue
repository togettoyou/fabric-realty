<template>
  <router-view></router-view>
</template>

<style>
#app {
  min-height: 100vh;
  background-color: #f0f2f5;
}

/* 内容区域样式 */
.content {
  margin: 0 auto;
  margin-top: 96px;
  max-width: 1200px;
  padding: 24px;
}

/* 全局卡片样式 */
.ant-card {
  border-radius: 8px;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.03);
}

/* 卡片内容区域样式 */
.ant-card-body {
  padding: 0 !important;
}

/* 卡片额外内容样式 */
.ant-card-extra {
  padding: 24px !important;
  border-bottom: 1px solid #f0f0f0;
}

/* 表格内的价格样式 */
.price {
  font-variant-numeric: tabular-nums;
  font-weight: 500;
}

</style>
