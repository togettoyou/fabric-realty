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

/* 全局表格样式 */
.ant-table {
  border-radius: 8px;
  overflow: hidden;
}

.ant-table-thead > tr > th {
  background: #fafafa !important;
  font-weight: 600 !important;
  color: #1f1f1f !important;
  border-bottom: 1px solid #f0f0f0 !important;
  padding: 16px !important;
}

.ant-table-tbody > tr > td {
  border-bottom: 1px solid #f0f0f0 !important;
  padding: 16px !important;
}

.ant-table-tbody > tr:hover > td {
  background: #fafafa !important;
}

.ant-table-row {
  transition: all 0.3s ease;
}

/* 表格内的标签样式 */
.ant-tag {
  border-radius: 4px;
  padding: 2px 8px;
  font-size: 13px;
}

/* 表格内的价格样式 */
.price {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto,
    'Helvetica Neue', Arial, 'Noto Sans', sans-serif, 'Apple Color Emoji',
    'Segoe UI Emoji', 'Segoe UI Symbol', 'Noto Color Emoji';
  font-variant-numeric: tabular-nums;
  font-weight: 500;
  color: #1f1f1f;
}

/* 表格内的时间样式 */
.ant-table-cell time {
  color: #666;
  font-size: 13px;
}

/* 表格内的按钮样式 */
.ant-table .ant-btn {
  border-radius: 4px;
}

/* 表格底部样式 */
.table-footer {
  padding: 16px;
  text-align: center;
  background: #fff;
  border-top: 1px solid #f0f0f0;
  border-radius: 0 0 8px 8px;
}

/* 全局卡片样式 */
.ant-card {
  border-radius: 8px;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.03), 0 1px 6px -1px rgba(0, 0, 0, 0.02), 0 2px 4px 0 rgba(0, 0, 0, 0.02);
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

</style>
