<template>
  <div class="realty-agency">
    <a-page-header
      title="不动产登记机构"
      sub-title="负责房产信息的登记和所有权变更"
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

    <div class="content">
      <a-card :bordered="false">
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
          </template>
        </a-table>
        <div class="table-footer">
          <a-button
            :disabled="!bookmark"
            @click="loadMore"
            :loading="loading"
            v-if="realEstateList.length > 0"
          >
            <template #icon><DownOutlined /></template>
            加载更多
          </a-button>
          <a-empty v-else />
        </div>
      </a-card>
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
          <a-input v-model:value="formState.address" placeholder="例如: 北京市朝阳区xxx街道xxx号" />
        </a-form-item>

        <a-form-item label="面积（平方米）" name="area" extra="请输入大于0的数值">
          <a-input-number
            v-model:value="formState.area"
            :min="0.01"
            :step="0.01"
            style="width: 100%"
            placeholder="请输入面积"
          />
        </a-form-item>

        <a-form-item label="所有者" name="owner" extra="可以输入任意模拟用户名">
          <a-input v-model:value="formState.owner" placeholder="请输入所有者姓名" />
        </a-form-item>

        <div class="form-tips">
          <InfoCircleOutlined style="color: #1890ff; margin-right: 8px;" />
          <span>房产ID将由系统自动生成</span>
        </div>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { message } from 'ant-design-vue';
import { PlusOutlined, EyeOutlined, DownOutlined, InfoCircleOutlined } from '@ant-design/icons-vue';
import { realtyApi } from '../api';
import type { FormInstance } from 'ant-design-vue';

const formRef = ref<FormInstance>();
const showCreateModal = ref(false);
const modalLoading = ref(false);

const formState = reactive({
  id: '',
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
  },
  {
    title: '地址',
    dataIndex: 'propertyAddress',
    key: 'propertyAddress',
    ellipsis: true,
  },
  {
    title: '面积（平方米）',
    dataIndex: 'area',
    key: 'area',
    width: 120,
    align: 'right' as const,
  },
  {
    title: '当前所有者',
    dataIndex: 'currentOwner',
    key: 'currentOwner',
    width: 120,
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

// 生成UUID
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
      // 自动生成房产ID
      const realEstateData = {
        ...formState,
        id: generateUUID(),
      };
      await realtyApi.createRealEstate(realEstateData);
      message.success('房产信息登记成功');
      showCreateModal.value = false;
      // 重置表单
      formRef.value?.resetFields();
      // 刷新列表
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

:deep(.ant-card-body) {
  padding: 0;
}

:deep(.ant-table-thead > tr > th) {
  background: #fafafa;
}

.table-footer {
  padding: 16px;
  text-align: center;
  background: #fff;
  border-radius: 0 0 8px 8px;
}

:deep(.ant-form-item-label) {
  font-weight: 500;
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