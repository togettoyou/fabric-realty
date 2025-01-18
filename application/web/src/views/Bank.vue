<template>
  <div class="bank">
    <div class="page-header">
      <a-page-header
        title="银行"
        sub-title="负责管理交易资金"
        @back="() => $router.push('/')"
      />
    </div>

    <div class="content">
      <a-card :bordered="false">
        <template #extra>
          <a-radio-group v-model:value="statusFilter" button-style="solid">
            <a-radio-button value="">全部</a-radio-button>
            <a-radio-button value="PENDING">待完成</a-radio-button>
            <a-radio-button value="COMPLETED">已完成</a-radio-button>
          </a-radio-group>
        </template>

        <div class="table-container">
          <a-table
            :columns="columns"
            :data-source="filteredTransactionList"
            :loading="loading"
            :pagination="false"
            :scroll="{ x: 1500, y: 'calc(100vh - 350px)' }"
            row-key="id"
            class="custom-table"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'id'">
                <div class="id-cell">
                  <a-tooltip :title="record.id">
                    <span class="id-text">{{ record.id }}</span>
                  </a-tooltip>
                  <a-tooltip title="点击复制">
                    <copy-outlined
                      class="copy-icon"
                      @click.stop="handleCopy(record.id)"
                    />
                  </a-tooltip>
                </div>
              </template>
              <template v-else-if="column.key === 'realEstateId'">
                <div class="id-cell">
                  <a-tooltip :title="record.realEstateId">
                    <span class="id-text">{{ record.realEstateId }}</span>
                  </a-tooltip>
                  <a-tooltip title="点击复制">
                    <copy-outlined
                      class="copy-icon"
                      @click.stop="handleCopy(record.realEstateId)"
                    />
                  </a-tooltip>
                </div>
              </template>
              <template v-else-if="column.key === 'status'">
                <a-tag :color="getStatusColor(record.status)">
                  {{ getStatusText(record.status) }}
                </a-tag>
              </template>
              <template v-else-if="column.key === 'createTime'">
                <time>{{ new Date(record.createTime).toLocaleString() }}</time>
              </template>
              <template v-else-if="column.key === 'updateTime'">
                <time>{{ new Date(record.updateTime).toLocaleString() }}</time>
              </template>
              <template v-else-if="column.key === 'action'">
                <a-space>
                  <a-tooltip title="完成交易">
                    <a-button
                      type="primary"
                      :disabled="record.status !== 'PENDING'"
                      @click="completeTransaction(record.id)"
                      :loading="record.loading"
                      size="small"
                    >
                      <template #icon><check-outlined /></template>
                      完成交易
                    </a-button>
                  </a-tooltip>
                </a-space>
              </template>
            </template>
          </a-table>
          <div class="load-more">
            <a-button 
              :loading="loading" 
              @click="loadMore"
              :disabled="!bookmark"
            >
              {{ bookmark ? '加载更多' : '没有更多数据了' }}
            </a-button>
          </div>
        </div>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { message } from 'ant-design-vue';
import { CheckOutlined, DownOutlined, CopyOutlined } from '@ant-design/icons-vue';
import { transactionApi } from '../api';

const statusFilter = ref('');
const transactionList = ref<any[]>([]);
const loading = ref(false);
const bookmark = ref('');

const handleCopy = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text);
    message.success('已复制到剪贴板');
  } catch (err) {
    message.error('复制失败');
  }
};

const columns = [
  {
    title: '交易ID',
    dataIndex: 'id',
    key: 'id',
    width: 180,
    ellipsis: false,
    customCell: () => ({
      style: {
        whiteSpace: 'nowrap',
        overflow: 'hidden',
      }
    }),
  },
  {
    title: '房产ID',
    dataIndex: 'realEstateId',
    key: 'realEstateId',
    width: 180,
    ellipsis: false,
    customCell: () => ({
      style: {
        whiteSpace: 'nowrap',
        overflow: 'hidden',
      }
    }),
  },
  {
    title: '卖家',
    dataIndex: 'seller',
    key: 'seller',
    width: 100,
    ellipsis: true,
  },
  {
    title: '买家',
    dataIndex: 'buyer',
    key: 'buyer',
    width: 100,
    ellipsis: true,
  },
  {
    title: '价格',
    dataIndex: 'price',
    key: 'price',
    width: 120,
    align: 'right' as const,
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    width: 80,
  },
  {
    title: '创建时间',
    dataIndex: 'createTime',
    key: 'createTime',
    width: 160,
  },
  {
    title: '更新时间',
    dataIndex: 'updateTime',
    key: 'updateTime',
    width: 160,
  },
  {
    title: '操作',
    key: 'action',
    fixed: 'right',
    width: 120,
    align: 'center' as const,
  },
];

const filteredTransactionList = computed(() => {
  if (!statusFilter.value) {
    return transactionList.value;
  }
  return transactionList.value.filter(item => item.status === statusFilter.value);
});

const loadTransactionList = async () => {
  try {
    loading.value = true;
    const result = await transactionApi.getTransactionList({
      pageSize: 10,
      bookmark: bookmark.value,
    });
    const records = result.records.map((record: any) => ({
      ...record,
      loading: false,
    }));
    transactionList.value = [...transactionList.value, ...records];
    bookmark.value = result.bookmark;
  } catch (error: any) {
    message.error(error.message || '加载交易列表失败');
  } finally {
    loading.value = false;
  }
};

const loadMore = () => {
  loadTransactionList();
};

const getStatusColor = (status: string) => {
  switch (status) {
    case 'PENDING':
      return 'blue';
    case 'COMPLETED':
      return 'green';
    default:
      return 'default';
  }
};

const getStatusText = (status: string) => {
  switch (status) {
    case 'PENDING':
      return '待完成';
    case 'COMPLETED':
      return '已完成';
    default:
      return '未知';
  }
};

const completeTransaction = async (txId: string) => {
  const transaction = transactionList.value.find((t) => t.id === txId);
  if (!transaction) return;

  try {
    transaction.loading = true;
    await transactionApi.completeTransaction(txId);
    message.success('交易完成');
    // 刷新列表
    transactionList.value = [];
    bookmark.value = '';
    loadTransactionList();
  } catch (error: any) {
    message.error(error.message || '完成交易失败');
  } finally {
    transaction.loading = false;
  }
};

// 初始加载
onMounted(() => {
  loadTransactionList();
});
</script>

<style scoped>
.bank {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.page-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  padding: 16px 24px;
}

.content {
  flex: 1;
  margin-top: 72px;
  padding: 24px;
  overflow: hidden;
}

.id-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.id-text {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
}

:deep(.copy-icon) {
  color: rgba(0, 0, 0, 0.45);
  font-size: 14px;
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.2s;
  
  &:hover {
    color: #1890ff;
  }
}

:deep(.ant-table-cell:hover .copy-icon) {
  opacity: 1;
}

.table-container {
  height: calc(100vh - 200px);
  position: relative;
  display: flex;
  flex-direction: column;
}

:deep(.ant-table-wrapper) {
  flex: 1;
  overflow: hidden;
}

:deep(.ant-table) {
  height: calc(100vh - 300px);
}

:deep(.ant-table-body) {
  max-height: calc(100vh - 360px) !important;
}

/* 固定表头样式 */
:deep(.ant-table-header) {
  background: #fff;
}

:deep(.ant-table-header::-webkit-scrollbar) {
  display: none;
}

.load-more {
  text-align: center;
  margin-top: 16px;
  padding: 8px 0;
  background: #fff;
  border-top: 1px solid #f0f0f0;
}
</style> 