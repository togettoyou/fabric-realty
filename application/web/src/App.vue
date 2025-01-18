<template>
  <router-view></router-view>
</template>

<style>
#app {
  min-height: 100vh;
  background-color: #f5f5f5;
  color: rgba(0, 0, 0, 0.88);
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
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.03);
}

/* 卡片内容区域样式 */
.ant-card-body {
  padding: 24px !important;
}

/* 页头样式 */
.page-header {
  background: #fff;
  border-bottom: 1px solid rgba(5, 5, 5, 0.06);
}

/* 表单样式 */
.ant-form-item-label > label {
  color: rgba(0, 0, 0, 0.88);
  font-size: 14px;
}

/* 按钮样式 */
.ant-btn {
  border-radius: 6px;
  font-size: 14px;
  height: 32px;
  padding: 4px 15px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.ant-btn-primary {
  box-shadow: 0 2px 0 rgba(5, 145, 255, 0.1);
}

/* 表格样式调整 */
.ant-table {
  color: rgba(0, 0, 0, 0.88);
}

.ant-table-thead > tr > th {
  background: #fafafa;
  font-weight: 600;
}

/* 标签样式 */
.ant-tag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 24px;
  padding: 0 8px;
  font-size: 13px;
  border-radius: 4px;
}

/* 价格样式 */
.price {
  font-variant-numeric: tabular-nums;
  font-weight: 500;
}

/* 文字颜色 */
.text-secondary {
  color: rgba(0, 0, 0, 0.45);
}

/* 滚动条美化 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: #f5f5f5;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb {
  background: #ccc;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: #999;
}
</style>
