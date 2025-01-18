<template>
  <div class="bank">
    <div class="page-header">
      <a-page-header
        title="银行"
        sub-title="负责交易的完成确认"
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

        <a-table
          :columns="columns"
          :data-source="filteredTransactionList"
          :loading="loading"
          :pagination="false"
          row-key="id"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'status'">
              <a-tag :color="getStatusColor(record.status)">
                {{ getStatusText(record.status) }}
              </a-tag>
            </template>
            <template v-else-if="column.key === 'price'">
              <span class="price">¥ {{ record.price.toLocaleString() }}</span>
            </template>
            <template v-else-if="column.key === 'createTime'">
              <time>{{ new Date(record.createTime).toLocaleString() }}</time>
            </template>
            <template v-else-if="column.key === 'updateTime'">
              <time>{{ new Date(record.updateTime).toLocaleString() }}</time>
            </template>
            <template v-else-if="column.key === 'action'">
              <a-space>
                <a-tooltip title="查看详情">
                  <a-button type="link" @click="viewDetail(record)">
                    <template #icon><EyeOutlined /></template>
                  </a-button>
                </a-tooltip>
                <a-tooltip title="完成交易">
                  <a-button
                    type="primary"
                    :disabled="record.status !== 'PENDING'"
                    @click="completeTransaction(record.id)"
                    :loading="record.loading"
                    size="small"
                  >
                    <template #icon><CheckOutlined /></template>
                    完成交易
                  </a-button>
                </a-tooltip>
              </a-space>
            </template>
          </template>
        </a-table>
        <div class="table-footer">
          <a-button
            :disabled="!bookmark"
            @click="loadMore"
            :loading="loading"
            v-if="transactionList.length > 0"
          >
            <template #icon><DownOutlined /></template>
            加载更多
          </a-button>
          <a-empty v-else />
        </div>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { message } from 'ant-design-vue';
import { EyeOutlined, CheckOutlined, DownOutlined } from '@ant-design/icons-vue';
import { transactionApi } from '../api';

const statusFilter = ref('');
const transactionList = ref<any[]>([]);
const loading = ref(false);
const bookmark = ref('');

const columns = [
  {
    title: '交易ID',
    dataIndex: 'id',
    key: 'id',
    width: 180,
    ellipsis: true,
    customCell: () => ({
      style: { cursor: 'copy' },
      onClick: (e: MouseEvent) => {
        const text = (e.target as HTMLElement).innerText;
        navigator.clipboard.writeText(text);
        message.success('已复制到剪贴板');
      },
    }),
  },
  {
    title: '房产ID',
    dataIndex: 'realEstateId',
    key: 'realEstateId',
    width: 180,
    ellipsis: true,
    customCell: () => ({
      style: { cursor: 'copy' },
      onClick: (e: MouseEvent) => {
        const text = (e.target as HTMLElement).innerText;
        navigator.clipboard.writeText(text);
        message.success('已复制到剪贴板');
      },
    }),
  },
  {
    title: '卖家',
    dataIndex: 'seller',
    key: 'seller',
    width: 120,
    ellipsis: true,
  },
  {
    title: '买家',
    dataIndex: 'buyer',
    key: 'buyer',
    width: 120,
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
    width: 100,
  },
  {
    title: '创建时间',
    dataIndex: 'createTime',
    key: 'createTime',
    width: 180,
  },
  {
    title: '更新时间',
    dataIndex: 'updateTime',
    key: 'updateTime',
    width: 180,
  },
  {
    title: '操作',
    key: 'action',
    width: 200,
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
    case 'CANCELLED':
      return 'red';
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
    case 'CANCELLED':
      return '已取消';
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

const viewDetail = (record: any) => {
  console.log('查看详情:', record);
  message.info('详情功能开发中');
};

// 初始加载
onMounted(() => {
  loadTransactionList();
});
</script>

<style scoped>
.bank {
  min-height: 100vh;
  background-color: #f0f2f5;
}

.page-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  background: #fff;
  padding: 16px 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.price {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto,
    'Helvetica Neue', Arial, 'Noto Sans', sans-serif, 'Apple Color Emoji',
    'Segoe UI Emoji', 'Segoe UI Symbol', 'Noto Color Emoji';
  font-variant-numeric: tabular-nums;
}
</style> 