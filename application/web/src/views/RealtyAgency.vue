<template>
  <div class="realty-agency">
    <a-page-header
      title="不动产登记机构"
      sub-title="负责房产信息的登记和所有权变更"
      @back="() => $router.push('/')"
    />

    <div class="content">
      <a-row :gutter="[24, 24]">
        <a-col :span="24">
          <a-card title="登记房产信息">
            <a-form
              :model="formState"
              name="createRealEstate"
              @finish="onFinish"
              @finishFailed="onFinishFailed"
            >
              <a-form-item
                label="房产ID"
                name="id"
                :rules="[{ required: true, message: '请输入房产ID' }]"
              >
                <a-input v-model:value="formState.id" />
              </a-form-item>

              <a-form-item
                label="房产地址"
                name="address"
                :rules="[{ required: true, message: '请输入房产地址' }]"
              >
                <a-input v-model:value="formState.address" />
              </a-form-item>

              <a-form-item
                label="面积"
                name="area"
                :rules="[{ required: true, message: '请输入面积' }]"
              >
                <a-input-number
                  v-model:value="formState.area"
                  :min="0"
                  :step="0.01"
                  style="width: 100%"
                />
              </a-form-item>

              <a-form-item
                label="所有者"
                name="owner"
                :rules="[{ required: true, message: '请输入所有者' }]"
              >
                <a-input v-model:value="formState.owner" />
              </a-form-item>

              <a-form-item>
                <a-button type="primary" html-type="submit">登记</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>

        <a-col :span="24">
          <a-card title="房产列表">
            <a-table
              :columns="columns"
              :data-source="realEstateList"
              :loading="loading"
              :pagination="false"
              row-key="id"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'status'">
                  <a-tag :color="record.status === 'NORMAL' ? 'green' : 'blue'">
                    {{ record.status === 'NORMAL' ? '正常' : '交易中' }}
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
import { realtyApi } from '../api';

const formState = reactive({
  id: '',
  address: '',
  area: 0,
  owner: '',
});

const columns = [
  {
    title: '房产ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '地址',
    dataIndex: 'propertyAddress',
    key: 'propertyAddress',
  },
  {
    title: '面积（平方米）',
    dataIndex: 'area',
    key: 'area',
  },
  {
    title: '当前所有者',
    dataIndex: 'currentOwner',
    key: 'currentOwner',
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

const realEstateList = ref<any[]>([]);
const loading = ref(false);
const bookmark = ref('');

const loadRealEstateList = async () => {
  try {
    loading.value = true;
    const result = await realtyApi.getRealEstateList({
      pageSize: 10,
      bookmark: bookmark.value,
    });
    realEstateList.value = [...realEstateList.value, ...result.records];
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

const onFinish = async (values: any) => {
  try {
    await realtyApi.createRealEstate(values);
    message.success('房产信息登记成功');
    // 重置表单
    formState.id = '';
    formState.address = '';
    formState.area = 0;
    formState.owner = '';
    // 刷新列表
    realEstateList.value = [];
    bookmark.value = '';
    loadRealEstateList();
  } catch (error: any) {
    message.error(error.message || '房产信息登记失败');
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
  loadRealEstateList();
});
</script>

<style scoped>
.realty-agency {
  padding: 24px;
  background-color: #f0f2f5;
  min-height: 100vh;
}

.content {
  margin-top: 24px;
}
</style> 