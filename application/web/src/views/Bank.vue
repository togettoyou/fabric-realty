<template>
  <div class="bank">
    <div class="app-page-header">
      <a-page-header
        title="银行"
        sub-title="负责交易的完成确认"
        @back="() => $router.push('/')"
      />
    </div>

    <div class="app-content">
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
              <template v-else-if="column.key === 'action'">
                <a-button
                  type="primary"
                  size="small"
                  @click="handleComplete(record)"
                  :loading="record.completing"
                  :disabled="record.status === 'COMPLETED'"
                >
                  完成交易
                </a-button>
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
import { CopyOutlined, ApartmentOutlined } from '@ant-design/icons-vue';
import { bankApi } from '../api';
import { ref, reactive } from 'vue';
import type { BlockData } from '../types';
import { copyToClipboard, getStatusText, getStatusColor, formatPrice } from '../utils';

const transactionList = ref<any[]>([]);
const loading = ref(false);
const bookmark = ref('');
const statusFilter = ref('');
const searchId = ref('');

const blockDrawer = ref(false);
const blockList = ref<BlockData[]>([]);
const blockTotal = ref(0);
const blockQuery = reactive({
  pageSize: 10,
  pageNum: 1,
});

const loadTransactionList = async () => {
  try {
    loading.value = true;
    const result = await bankApi.getTransactionList({
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

const handleComplete = async (record: any) => {
  try {
    record.completing = true;
    await bankApi.completeTransaction(record.id);
    message.success('交易完成');
    // 刷新列表
    transactionList.value = [];
    bookmark.value = '';
    loadTransactionList();
  } catch (error: any) {
    message.error(error.message || '完成交易失败');
  } finally {
    record.completing = false;
  }
};

const handleCopy = (text: string) => {
  copyToClipboard(text);
};

const handleSearch = async (value: string) => {
  if (!value) {
    message.warning('请输入要查询的交易ID');
    return;
  }

  try {
    const result = await bankApi.getTransaction(value);
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

// 监听状态筛选变化
watch(statusFilter, () => {
  // 重置列表和书签
  transactionList.value = [];
  bookmark.value = '';
  // 重新加载数据
  loadTransactionList();
});

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
    customRender: ({ text }: { text: number }) => formatPrice(text),
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
    width: 100,
    fixed: 'right',
  },
];

onMounted(() => {
  loadTransactionList();
});

// 打开区块抽屉
const openBlockDrawer = async () => {
  blockDrawer.value = true;
  await fetchBlockList();
};

// 获取区块列表
const fetchBlockList = async () => {
  try {
    const res = await bankApi.getBlockList({
      pageSize: blockQuery.pageSize,
      pageNum: blockQuery.pageNum,
    });
    blockList.value = res.blocks;
    blockTotal.value = res.total;
  } catch (error) {
    console.error('获取区块列表失败：', error);
  }
};

// 处理分页变化
const handleBlockPageChange = async (page: number, pageSize: number) => {
  blockQuery.pageNum = page;
  blockQuery.pageSize = pageSize;
  await fetchBlockList();
};
</script>

<style scoped>
</style>
