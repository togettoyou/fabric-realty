<template>
  <div class="trading-platform">
    <div class="page-header">
      <a-page-header
        title="交易平台"
        sub-title="负责创建和管理交易信息"
        @back="() => $router.push('/')"
      >
        <template #extra>
          <a-tooltip title="点击创建新的交易">
            <a-button type="primary" @click="showCreateModal = true">
              <template #icon><PlusOutlined /></template>
              创建新交易
            </a-button>
          </a-tooltip>
        </template>
      </a-page-header>
    </div>

    <div class="content">
      <a-card :bordered="false">
        <a-table
          :columns="columns"
          :data-source="transactionList"
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

    <!-- 创建交易的对话框 -->
    <a-modal
      v-model:visible="showCreateModal"
      title="创建新交易"
      @ok="handleModalOk"
      @cancel="handleModalCancel"
      :confirmLoading="modalLoading"
    >
      <a-form
        ref="formRef"
        :model="formState"
        :rules="rules"
        layout="vertical"
      >
        <a-form-item label="房产ID" name="realEstateId" extra="请输入要交易的房产ID">
          <a-input v-model:value="formState.realEstateId" placeholder="请输入房产ID" />
        </a-form-item>

        <a-form-item label="卖家" name="seller" extra="可以输入任意模拟用户名作为卖家">
          <a-input v-model:value="formState.seller" placeholder="请输入卖家姓名" />
        </a-form-item>

        <a-form-item label="买家" name="buyer" extra="可以输入任意模拟用户名作为买家">
          <a-input v-model:value="formState.buyer" placeholder="请输入买家姓名" />
        </a-form-item>

        <a-form-item label="价格" name="price" extra="请输入大于0的交易金额">
          <a-input-number
            v-model:value="formState.price"
            :min="0.01"
            :step="0.01"
            style="width: 100%"
            placeholder="请输入价格"
            :formatter="value => `¥ ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')"
            :parser="value => value!.replace(/\¥\s?|(,*)/g, '')"
          />
        </a-form-item>

        <div class="form-tips">
          <InfoCircleOutlined style="color: #1890ff; margin-right: 8px;" />
          <span>交易ID将由系统自动生成</span>
        </div>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { message } from 'ant-design-vue';
import { PlusOutlined, EyeOutlined, DownOutlined, InfoCircleOutlined } from '@ant-design/icons-vue';
import { transactionApi } from '../api';
import type { FormInstance } from 'ant-design-vue';

const formRef = ref<FormInstance>();
const showCreateModal = ref(false);
const modalLoading = ref(false);

const formState = reactive({
  txId: '',
  realEstateId: '',
  seller: '',
  buyer: '',
  price: undefined as number | undefined,
});

const rules = {
  realEstateId: [{ required: true, message: '请输入房产ID' }],
  seller: [{ required: true, message: '请输入卖家' }],
  buyer: [{ required: true, message: '请输入买家' }],
  price: [
    { required: true, message: '请输入价格' },
    { type: 'number', min: 0.01, message: '价格必须大于0' }
  ],
};

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
];

const transactionList = ref<any[]>([]);
const loading = ref(false);
const bookmark = ref('');

const loadTransactionList = async () => {
  try {
    loading.value = true;
    const result = await transactionApi.getTransactionList({
      pageSize: 10,
      bookmark: bookmark.value,
    });
    transactionList.value = [...transactionList.value, ...result.records];
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

const generateUUID = () => {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    const r = Math.random() * 16 | 0;
    const v = c === 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
};

const handleModalOk = () => {
  formRef.value?.validate().then(async () => {
    modalLoading.value = true;
    try {
      const transactionData = {
        ...formState,
        txId: generateUUID(),
      };
      await transactionApi.createTransaction(transactionData);
      message.success('交易创建成功');
      showCreateModal.value = false;
      formRef.value?.resetFields();
      transactionList.value = [];
      bookmark.value = '';
      loadTransactionList();
    } catch (error: any) {
      message.error(error.message || '交易创建失败');
    } finally {
      modalLoading.value = false;
    }
  });
};

const handleModalCancel = () => {
  showCreateModal.value = false;
  formRef.value?.resetFields();
};

// 初始加载
onMounted(() => {
  loadTransactionList();
});
</script>

<style scoped>
.trading-platform {
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

:deep(.ant-form-item-label) {
  font-weight: 500;
}

.price {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto,
    'Helvetica Neue', Arial, 'Noto Sans', sans-serif, 'Apple Color Emoji',
    'Segoe UI Emoji', 'Segoe UI Symbol', 'Noto Color Emoji';
  font-variant-numeric: tabular-nums;
}

.form-tips {
  background-color: #e6f7ff;
  padding: 8px 12px;
  border-radius: 4px;
  color: #666;
  font-size: 14px;
  display: flex;
  align-items: center;
}

:deep(.ant-form-item-extra) {
  color: #666;
}
</style> 