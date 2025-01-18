<template>
  <div class="trading-platform">
    <a-page-header
      title="交易平台"
      sub-title="负责创建和管理交易信息"
      @back="() => $router.push('/')"
    />

    <div class="content">
      <a-row :gutter="[24, 24]">
        <a-col :span="24">
          <a-card title="创建交易">
            <a-form
              :model="formState"
              name="createTransaction"
              @finish="onFinish"
              @finishFailed="onFinishFailed"
            >
              <a-form-item
                label="交易ID"
                name="txId"
                :rules="[{ required: true, message: '请输入交易ID' }]"
              >
                <a-input v-model:value="formState.txId" />
              </a-form-item>

              <a-form-item
                label="房产ID"
                name="realEstateId"
                :rules="[{ required: true, message: '请输入房产ID' }]"
              >
                <a-input v-model:value="formState.realEstateId" />
              </a-form-item>

              <a-form-item
                label="卖家"
                name="seller"
                :rules="[{ required: true, message: '请输入卖家' }]"
              >
                <a-input v-model:value="formState.seller" />
              </a-form-item>

              <a-form-item
                label="买家"
                name="buyer"
                :rules="[{ required: true, message: '请输入买家' }]"
              >
                <a-input v-model:value="formState.buyer" />
              </a-form-item>

              <a-form-item
                label="价格"
                name="price"
                :rules="[{ required: true, message: '请输入价格' }]"
              >
                <a-input-number
                  v-model:value="formState.price"
                  :min="0"
                  :step="0.01"
                  style="width: 100%"
                />
              </a-form-item>

              <a-form-item>
                <a-button type="primary" html-type="submit">创建</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>

        <a-col :span="24">
          <a-card title="交易列表">
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
                <template v-else-if="column.key === 'createTime'">
                  {{ new Date(record.createTime).toLocaleString() }}
                </template>
                <template v-else-if="column.key === 'updateTime'">
                  {{ new Date(record.updateTime).toLocaleString() }}
                </template>
                <template v-else-if="column.key === 'action'">
                  <a-button type="link" @click="viewDetail(record)">查看</a-button>
                </template>
              </template>
            </a-table>
            <div style="margin-top: 16px; text-align: right">
              <a-button
                :disabled="!bookmark"
                @click="loadMore"
                :loading="loading"
              >
                加载更多
              </a-button>
            </div>
          </a-card>
        </a-col>
      </a-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { message } from 'ant-design-vue';
import { transactionApi } from '../api';

const formState = reactive({
  txId: '',
  realEstateId: '',
  seller: '',
  buyer: '',
  price: 0,
});

const columns = [
  {
    title: '交易ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '房产ID',
    dataIndex: 'realEstateId',
    key: 'realEstateId',
  },
  {
    title: '卖家',
    dataIndex: 'seller',
    key: 'seller',
  },
  {
    title: '买家',
    dataIndex: 'buyer',
    key: 'buyer',
  },
  {
    title: '价格',
    dataIndex: 'price',
    key: 'price',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
  },
  {
    title: '创建时间',
    dataIndex: 'createTime',
    key: 'createTime',
  },
  {
    title: '更新时间',
    dataIndex: 'updateTime',
    key: 'updateTime',
  },
  {
    title: '操作',
    key: 'action',
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

const onFinish = async (values: any) => {
  try {
    await transactionApi.createTransaction(values);
    message.success('交易创建成功');
    // 重置表单
    formState.txId = '';
    formState.realEstateId = '';
    formState.seller = '';
    formState.buyer = '';
    formState.price = 0;
    // 刷新列表
    transactionList.value = [];
    bookmark.value = '';
    loadTransactionList();
  } catch (error: any) {
    message.error(error.message || '交易创建失败');
  }
};

const onFinishFailed = (errorInfo: any) => {
  console.log('Failed:', errorInfo);
};

const viewDetail = (record: any) => {
  console.log('查看详情:', record);
};

// 初始加载
onMounted(() => {
  loadTransactionList();
});
</script>

<style scoped>
.trading-platform {
  padding: 24px;
  background-color: #f0f2f5;
  min-height: 100vh;
}

.content {
  margin-top: 24px;
}
</style> 