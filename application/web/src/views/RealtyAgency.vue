<template>
  <div class="realty-agency">
    <div class="app-page-header">
      <a-page-header
        title="不动产登记机构"
        sub-title="负责房产信息的登记"
        @back="() => $router.push('/')"
      >
        <template #extra>
          <a-tooltip title="点击创建新的房产信息">
            <a-button type="primary" @click="showCreateModal = true">
              <template #icon><PlusOutlined /></template>
              登记新房产
            </a-button>
          </a-tooltip>
        </template>
      </a-page-header>
    </div>

    <div class="app-content">
      <a-card :bordered="false">
        <template #extra>
          <div class="card-extra">
            <a-input-search
              v-model:value="searchId"
              placeholder="输入房产ID进行精确查询"
              style="width: 300px; margin-right: 16px;"
              @search="handleSearch"
              @change="handleSearchChange"
              allow-clear
            />
            <a-radio-group v-model:value="statusFilter" button-style="solid">
              <a-radio-button value="">全部</a-radio-button>
              <a-radio-button value="NORMAL">正常</a-radio-button>
              <a-radio-button value="IN_TRANSACTION">交易中</a-radio-button>
            </a-radio-group>
          </div>
        </template>

        <div class="table-container">
          <a-table
            :columns="columns"
            :data-source="realEstateList"
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
              <template v-else-if="column.key === 'currentOwner'">
                <div class="id-cell">
                  <a-tooltip :title="record.currentOwner">
                    <span class="id-text">{{ record.currentOwner }}</span>
                  </a-tooltip>
                  <a-tooltip title="点击复制">
                    <copy-outlined
                      class="copy-icon"
                      @click.stop="handleCopy(record.currentOwner)"
                    />
                  </a-tooltip>
                </div>
              </template>
              <template v-else-if="column.key === 'status'">
                <a-tag :color="record.status === 'NORMAL' ? 'green' : 'blue'">
                  {{ record.status === 'NORMAL' ? '正常' : '交易中' }}
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

    <div
      class="block-icon"
      @click="openBlockDrawer"
    >
      <ApartmentOutlined />
    </div>

    <!-- 创建房产的对话框 -->
    <a-modal
      v-model:visible="showCreateModal"
      title="登记新房产"
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
        <a-form-item label="房产地址" name="address" extra="请输入完整的房产地址信息">
          <a-input-group compact>
            <a-input
              v-model:value="formState.address"
              placeholder="例如: 北京市朝阳区xxx街道xxx号"
              style="width: calc(100% - 110px)"
            />
            <a-tooltip title="随机生成一个地址">
              <a-button @click="generateRandomAddressHandler">
                <template #icon><ReloadOutlined /></template>
                随机生成
              </a-button>
            </a-tooltip>
          </a-input-group>
        </a-form-item>

        <a-form-item label="面积（平方米）" name="area" extra="请输入大于0的数值">
          <a-input-group compact>
            <a-input-number
              v-model:value="formState.area"
              :min="0.01"
              :step="0.01"
              style="width: calc(100% - 110px)"
              placeholder="请输入面积"
            />
            <a-tooltip title="随机生成一个面积">
              <a-button @click="generateRandomAreaHandler">
                <template #icon><ReloadOutlined /></template>
                随机生成
              </a-button>
            </a-tooltip>
          </a-input-group>
        </a-form-item>

        <a-form-item label="所有者" name="owner" extra="可以输入任意模拟用户名">
          <a-input-group compact>
            <a-input
              v-model:value="formState.owner"
              placeholder="请输入所有者姓名"
              style="width: calc(100% - 110px)"
            />
            <a-tooltip title="随机生成一个用户名">
              <a-button @click="generateRandomOwner">
                <template #icon><ReloadOutlined /></template>
                随机生成
              </a-button>
            </a-tooltip>
          </a-input-group>
        </a-form-item>

        <div class="form-tips">
          <InfoCircleOutlined style="color: #1890ff; margin-right: 8px;" />
          <span>房产ID将由系统自动生成</span>
        </div>
      </a-form>
    </a-modal>

    <!-- 区块信息抽屉 -->
    <a-drawer
      v-model:visible="blockDrawer"
      title="区块信息"
      placement="right"
      :width="960"
      :closable="true"
    >
      <div class="block-container">
        <div class="block-header">
          <h3>区块列表</h3>
          <a-pagination
            v-model:current="blockQuery.pageNum"
            v-model:pageSize="blockQuery.pageSize"
            :total="blockTotal"
            :show-total="total => `共 ${total} 条`"
            :page-size-options="['5', '10', '20', '50']"
            show-size-changer
            @change="handleBlockPageChange"
          />
        </div>

        <div class="block-list">
          <a-card v-for="block in blockList" :key="block.block_num" class="block-item">
            <template #title>
              <div class="block-item-header">
                <span class="block-number">区块 #{{ block.block_num }}</span>
                <span class="block-time">{{ new Date(block.save_time).toLocaleString() }}</span>
              </div>
            </template>
            <div class="block-item-content">
              <div class="block-field">
                <span class="field-label">区块哈希：</span>
                <a-tooltip :title="block.block_hash">
                  <span class="field-value hash">{{ block.block_hash }}</span>
                </a-tooltip>
                <a-tooltip title="点击复制">
                  <copy-outlined
                    class="copy-icon"
                    @click="handleCopy(block.block_hash)"
                  />
                </a-tooltip>
              </div>
              <div class="block-field">
                <span class="field-label">数据哈希：</span>
                <a-tooltip :title="block.data_hash">
                  <span class="field-value hash">{{ block.data_hash }}</span>
                </a-tooltip>
                <a-tooltip title="点击复制">
                  <copy-outlined
                    class="copy-icon"
                    @click="handleCopy(block.data_hash)"
                  />
                </a-tooltip>
              </div>
              <div class="block-field">
                <span class="field-label">前块哈希：</span>
                <a-tooltip :title="block.prev_hash">
                  <span class="field-value hash">{{ block.prev_hash }}</span>
                </a-tooltip>
                <a-tooltip title="点击复制">
                  <copy-outlined
                    class="copy-icon"
                    @click="handleCopy(block.prev_hash)"
                  />
                </a-tooltip>
              </div>
              <div class="block-field">
                <span class="field-label">交易数量：</span>
                <span class="field-value">{{ block.tx_count }}</span>
              </div>
            </div>
          </a-card>
        </div>
      </div>
    </a-drawer>

  </div>
</template>

<script setup lang="ts">
import { message } from 'ant-design-vue';
import { PlusOutlined, InfoCircleOutlined, ReloadOutlined, CopyOutlined, ApartmentOutlined } from '@ant-design/icons-vue';
import { realtyAgencyApi } from '../api';
import type { FormInstance } from 'ant-design-vue';
import { ref, reactive } from 'vue';
import type { BlockData } from '../types';
import { copyToClipboard, generateRandomName, generateRandomAddress, generateRandomArea, generateUUID } from '../utils';

const formRef = ref<FormInstance>();
const showCreateModal = ref(false);
const modalLoading = ref(false);

const formState = reactive({
  address: '',
  area: undefined as number | undefined,
  owner: '',
});

const rules = {
  address: [{ required: true, message: '请输入房产地址' }],
  area: [
    { required: true, message: '请输入面积' },
    { type: 'number', min: 0.01, message: '面积必须大于0' }
  ],
  owner: [{ required: true, message: '请输入所有者' }],
};

const columns = [
  {
    title: '房产ID',
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
    title: '地址',
    dataIndex: 'propertyAddress',
    key: 'propertyAddress',
    width: 120,
    ellipsis: { showTitle: true },
  },
  {
    title: '面积',
    dataIndex: 'area',
    key: 'area',
    width: 80,
    customCell: () => ({
      style: {
        fontVariantNumeric: 'tabular-nums',
      },
    }),
    customRender: ({ text }: { text: number }) => `${text} ㎡`,
  },
  {
    title: '当前所有者',
    dataIndex: 'currentOwner',
    key: 'currentOwner',
    width: 120,
    ellipsis: false,
    customCell: () => ({
      style: {
        whiteSpace: 'nowrap',
        overflow: 'hidden',
      }
    }),
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

const realEstateList = ref<any[]>([]);
const loading = ref(false);
const bookmark = ref('');

const loadRealEstateList = async () => {
  try {
    loading.value = true;
    const result = await realtyAgencyApi.getRealEstateList({
      pageSize: 10,
      bookmark: bookmark.value,
      status: statusFilter.value,
    });
    if (!bookmark.value) {
      realEstateList.value = result.records;
    } else {
      realEstateList.value = [...realEstateList.value, ...result.records];
    }
    bookmark.value = result.bookmark;
  } catch (error: any) {
    message.error(error.message || '加载房产列表失败');
  } finally {
    loading.value = false;
  }
};

const loadMore = () => {
  loadRealEstateList();
};

const handleModalOk = () => {
  formRef.value?.validate().then(async () => {
    modalLoading.value = true;
    try {
      const realEstateData = {
        ...formState,
        id: generateUUID(),
      };
      await realtyAgencyApi.createRealEstate(realEstateData);
      message.success('房产信息登记成功');
      showCreateModal.value = false;
      formRef.value?.resetFields();
      realEstateList.value = [];
      bookmark.value = '';
      loadRealEstateList();
    } catch (error: any) {
      message.error(error.message || '房产信息登记失败');
    } finally {
      modalLoading.value = false;
    }
  });
};

const handleModalCancel = () => {
  showCreateModal.value = false;
  formRef.value?.resetFields();
};

const generateRandomAddressHandler = () => {
  formState.address = generateRandomAddress();
};

const generateRandomAreaHandler = () => {
  formState.area = generateRandomArea();
};

const generateRandomOwner = () => {
  formState.owner = generateRandomName();
};

const handleCopy = (text: string) => {
  copyToClipboard(text);
};

const statusFilter = ref('');

watch(statusFilter, () => {
  realEstateList.value = [];
  bookmark.value = '';
  loadRealEstateList();
});

const searchId = ref('');

const handleSearch = async (value: string) => {
  if (!value) {
    message.warning('请输入要查询的房产ID');
    return;
  }

  try {
    const result = await realtyAgencyApi.getRealEstate(value);
    realEstateList.value = [result];
    bookmark.value = '';
  } catch (error: any) {
    message.error(error.message || '查询房产信息失败');
    realEstateList.value = [];
  }
};

const handleSearchChange = (e: Event) => {
  const value = (e.target as HTMLInputElement).value;
  if (!value) {
    realEstateList.value = [];
    bookmark.value = '';
    loadRealEstateList();
  }
};

onMounted(() => {
  loadRealEstateList();
});

const blockDrawer = ref(false);
const blockList = ref<BlockData[]>([]);
const blockTotal = ref(0);
const blockQuery = reactive({
  pageSize: 10,
  pageNum: 1,
});

const openBlockDrawer = async () => {
  blockDrawer.value = true;
  await fetchBlockList();
};

const fetchBlockList = async () => {
  try {
    const res = await realtyAgencyApi.getBlockList({
      pageSize: blockQuery.pageSize,
      pageNum: blockQuery.pageNum,
    });
    blockList.value = res.blocks;
    blockTotal.value = res.total;
  } catch (error) {
    console.error('获取区块列表失败：', error);
  }
};

const handleBlockPageChange = async (page: number, pageSize: number) => {
  blockQuery.pageNum = page;
  blockQuery.pageSize = pageSize;
  await fetchBlockList();
};

</script>

<style scoped>
</style>
