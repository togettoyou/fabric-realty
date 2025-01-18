<template>
  <router-view></router-view>
</template>

<style>
#app {
  min-height: 100vh;
  background-color: #fff;
  color: rgba(0, 0, 0, 0.88);
  overflow: hidden;
}

/* 内容区域样式 */
.content {
  margin-top: 72px;
  padding: 24px;
}

/* 全局卡片样式 */
.ant-card {
  background: #fff;
  border-radius: 8px;
  box-shadow: none;
  border: none;
}

/* 卡片内容区域样式 */
.ant-card-body {
  padding: 24px !important;
}

/* 页头样式 */
.page-header {
  background: #fff;
  border-bottom: none;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.03);
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
  border-bottom: none !important;
}

.ant-table-tbody > tr > td {
  border-bottom: 1px solid #f5f5f5 !important;
}

.ant-table-tbody > tr:last-child > td {
  border-bottom: none !important;
}

/* 去掉表格外边框 */
.ant-table-container {
  border: none !important;
}

.ant-table-cell {
  border-right: none !important;
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
  border: none;
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

/* 表格滚动条样式 */
.table-container::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

.table-container::-webkit-scrollbar-track {
  background: #fff;
  border-radius: 3px;
}

.table-container::-webkit-scrollbar-thumb {
  background: #eee;
  border-radius: 3px;
}

.table-container::-webkit-scrollbar-thumb:hover {
  background: #ddd;
}

/* 单选按钮组样式 */
.ant-radio-group-solid .ant-radio-button-wrapper {
  border: none !important;
  background: #f5f5f5;
  margin-right: 8px;
}

.ant-radio-group-solid .ant-radio-button-wrapper:not(:first-child)::before {
  display: none;
}

.ant-radio-group-solid .ant-radio-button-wrapper-checked:not(.ant-radio-button-wrapper-disabled) {
  background: #1890ff;
}
</style>
