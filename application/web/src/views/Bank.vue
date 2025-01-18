<template>
  <div class="bank">
    <a-page-header
      title="银行"
      sub-title="负责交易的完成确认"
      @back="$router.push('/')"
    />
    <div class="content">
      <a-row :gutter="[24, 24]">
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
                  <a-tag :color="record.status === 'COMPLETED' ? 'green' : 'blue'">
                    {{ record.status === 'COMPLETED' ? '已完成' : '待付款' }}
                  </a-tag>
                </template>
                <template v-else-if="column.key === 'createTime'">
                  {{ new Date(record.createTime).toLocaleString() }}
                </template>
                <template v-else-if="column.key === 'updateTime'">
                  {{ new Date(record.updateTime).toLocaleString() }}
                </template>
                <template v-else-if="column.key === 'action'">
                  <a-button
                    type="primary"
                    size="small"
                    :disabled="record.status === 'COMPLETED'"
                    @click="completeTransaction(record.id)"
                    :loading="record.loading"
                  >
                    完成交易
                  </a-button>
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
import { message } from 'ant-design-vue'
import { transactionApi } from '../api'

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
]

const transactionList = ref<any[]>([])
const loading = ref(false)
const bookmark = ref('')

const loadTransactionList = async () => {
  try {
    loading.value = true
    const result = await transactionApi.getTransactionList({
      pageSize: 10,
      bookmark: bookmark.value,
    })
    const records = result.records.map((record: any) => ({
      ...record,
      loading: false,
    }))
    transactionList.value = [...transactionList.value, ...records]
    bookmark.value = result.bookmark
  } catch (error: any) {
    message.error(error.message || '加载交易列表失败')
  } finally {
    loading.value = false
  }
}

const loadMore = () => {
  loadTransactionList()
}

const completeTransaction = async (txId: string) => {
  const transaction = transactionList.value.find((t) => t.id === txId)
  if (!transaction) return

  try {
    transaction.loading = true
    await transactionApi.completeTransaction(txId)
    message.success('交易完成')
    // 刷新列表
    transactionList.value = []
    bookmark.value = ''
    loadTransactionList()
  } catch (error: any) {
    message.error(error.message || '完成交易失败')
  } finally {
    transaction.loading = false
  }
}

// 初始加载
onMounted(() => {
  loadTransactionList()
})
</script>

<style scoped>
.bank {
  padding: 24px;
  background-color: #f0f2f5;
  min-height: 100vh;
}

.content {
  margin-top: 24px;
}
</style> 