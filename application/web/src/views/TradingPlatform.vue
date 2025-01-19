<template>
  <div class="trading-platform">
    <div class="page-header">
      <a-page-header
        title="交易平台"
        sub-title="负责生成交易信息"
        @back="() => $router.push('/')"
      >
        <template #extra>
          <a-tooltip title="点击创建新的交易">
            <a-button type="primary" @click="showCreateModal = true">
              <template #icon><PlusOutlined /></template>
              生成新交易
            </a-button>
          </a-tooltip>
        </template>
      </a-page-header>
    </div>

    <div class="content">
      <a-card :bordered="false">
        <template #extra>
          <div class="card-extra">
            <a-input-search
              v-model:value="searchId"
              placeholder="输入交易ID进行精确查询"
              style="width: 300px; margin-right: 16px;"
              @search="handleSearch"
              @change="handleSearchChange"
              allow-clear
            />
            <a-radio-group v-model:value="statusFilter" button-style="solid">
              <a-radio-button value="">全部</a-radio-button>
              <a-radio-button value="PENDING">待完成</a-radio-button>
              <a-radio-button value="COMPLETED">已完成</a-radio-button>
            </a-radio-group>
          </div>
        </template>

        <div class="table-container">
          <a-table
            :columns="columns"
            :data-source="transactionList"
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

    <!-- 生成交易的对话框 -->
    <a-modal
      v-model:visible="showCreateModal"
      title="生成新交易"
      @ok="handleModalOk"
      @cancel="handleModalCancel"
      :confirmLoading="modalLoading"
      :style="{ top: '40px' }"
    >
      <a-form
        ref="formRef"
        :model="formState"
        :rules="rules"
        layout="vertical"
      >
        <a-form-item label="房产ID" name="realEstateId" extra="请输入要交易的房产ID">
          <a-input 
            v-model:value="formState.realEstateId" 
            placeholder="请输入房产ID" 
            @change="handleRealEstateIdChange"
          />
        </a-form-item>

        <a-form-item label="卖家" name="seller" extra="当前房产所有者">
          <a-input 
            v-model:value="formState.seller" 
            placeholder="自动填入当前所有者" 
            disabled
          />
        </a-form-item>

        <a-form-item label="买家" name="buyer" extra="可以输入任意模拟用户名作为买家">
          <a-input-group compact>
            <a-input 
              v-model:value="formState.buyer" 
              placeholder="请输入买家姓名"
              style="width: calc(100% - 110px)"
            />
            <a-tooltip title="随机生成一个用户名">
              <a-button @click="generateRandomBuyer">
                <template #icon><ReloadOutlined /></template>
                随机生成
              </a-button>
            </a-tooltip>
          </a-input-group>
        </a-form-item>

        <a-form-item label="价格" name="price" extra="请输入大于0的交易金额">
          <a-input-group compact>
            <a-input-number
              v-model:value="formState.price"
              :min="0.01"
              :step="0.01"
              style="width: calc(100% - 110px)"
              placeholder="请输入价格"
              :formatter="value => `¥ ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')"
              :parser="value => value!.replace(/\¥\s?|(,*)/g, '')"
            />
            <a-tooltip title="随机生成一个价格">
              <a-button @click="generateRandomPrice">
                <template #icon><ReloadOutlined /></template>
                随机生成
              </a-button>
            </a-tooltip>
          </a-input-group>
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
import { PlusOutlined, InfoCircleOutlined, CopyOutlined, ReloadOutlined, SearchOutlined } from '@ant-design/icons-vue';
import { transactionApi, realtyApi } from '../api';
import type { FormInstance } from 'ant-design-vue';

const formRef = ref<FormInstance>();
const showCreateModal = ref(false);
const modalLoading = ref(false);

const formState = reactive({
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
    customRender: ({ text }: { text: number }) => 
      `¥ ${text}`.replace(/\B(?=(\d{3})+(?!\d))/g, ','),
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
      status: statusFilter.value,
    });
    if (!bookmark.value) {
      transactionList.value = result.records;
    } else {
      transactionList.value = [...transactionList.value, ...result.records];
    }
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

const statusFilter = ref('');

watch(statusFilter, () => {
  transactionList.value = [];
  bookmark.value = '';
  loadTransactionList();
});

const handleCopy = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text);
    message.success('已复制到剪贴板');
  } catch (err) {
    message.error('复制失败');
  }
};

const handleRealEstateIdChange = async (e: Event) => {
  const id = (e.target as HTMLInputElement).value;
  if (!id) {
    formState.seller = '';
    return;
  }

  try {
    const result = await realtyApi.getRealEstate(id);
    if (result.status !== 'NORMAL') {
      message.error('该房产不是正常状态，无法生成交易');
      formState.realEstateId = '';
      formState.seller = '';
      return;
    }
    formState.seller = result.currentOwner;
  } catch (error: any) {
    message.error(error.message || '获取房产信息失败');
    formState.seller = '';
  }
};

// 随机生成买家姓名
const lastNames = ['张', '王', '李', '赵', '刘', '陈', '杨', '黄', '周', '吴'];
const firstNames = ['伟', '芳', '娜', '秀英', '敏', '静', '丽', '强', '磊', '洋', '艳', '勇', '军', '杰', '娟', '涛', '明', '超', '秀兰', '霞'];

const generateRandomBuyer = () => {
  const lastName = lastNames[Math.floor(Math.random() * lastNames.length)];
  const firstName = firstNames[Math.floor(Math.random() * firstNames.length)];
  formState.buyer = lastName + firstName;
};

// 随机生成价格
const generateRandomPrice = () => {
  // 生成 50-1000 万之间的随机价格（单位：元），保留2位小数
  formState.price = Number((Math.random() * (10000000 - 500000) + 500000).toFixed(2));
};

// 添加搜索相关的变量和方法
const searchId = ref('');

const handleSearch = async (value: string) => {
  if (!value) {
    message.warning('请输入要查询的交易ID');
    return;
  }
  
  try {
    const result = await transactionApi.getTransaction(value);
    transactionList.value = [result];
    bookmark.value = '';
  } catch (error: any) {
    message.error(error.message || '查询交易信息失败');
    transactionList.value = [];
  }
};

const handleSearchChange = (e: Event) => {
  const value = (e.target as HTMLInputElement).value;
  if (!value) {
    // 当搜索框清空时，重新加载列表
    transactionList.value = [];
    bookmark.value = '';
    loadTransactionList();
  }
};

onMounted(() => {
  loadTransactionList();
});
</script>

<style scoped>
.trading-platform {
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

:deep(.form-tips) {
  background-color: #e6f7ff;
  padding: 8px 12px;
  border-radius: 4px;
  color: rgba(0, 0, 0, 0.65);
  font-size: 14px;
  display: flex;
  align-items: center;
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

.load-more {
  text-align: center;
  margin-top: 16px;
  padding: 8px 0;
  background: #fff;
  border-top: 1px solid #f0f0f0;
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

:deep(.ant-table-header) {
  background: #fff;
}

:deep(.ant-table-header::-webkit-scrollbar) {
  display: none;
}

.card-extra {
  display: flex;
  align-items: center;
  gap: 16px;
}
</style> 