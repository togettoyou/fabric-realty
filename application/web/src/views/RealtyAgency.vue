<template>
  <div class="realty-agency">
    <div class="page-header">
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

    <div class="content">
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
              <a-button @click="generateRandomAddress">
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
              <a-button @click="generateRandomArea">
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
  </div>
</template>

<script setup lang="ts">
import { message } from 'ant-design-vue';
import { PlusOutlined, InfoCircleOutlined, ReloadOutlined, CopyOutlined, SearchOutlined } from '@ant-design/icons-vue';
import { realtyAgencyApi } from '../api';
import type { FormInstance } from 'ant-design-vue';

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
      await realtyAgencyApi.createRealEstate(realEstateData);
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

// 随机生成地址
const cities = ['北京市', '上海市', '广州市', '深圳市', '杭州市', '南京市', '成都市', '武汉市'];
const districts = ['东城区', '西城区', '朝阳区', '海淀区', '丰台区', '昌平区'];
const streets = ['长安街', '建国路', '复兴路', '中关村大街', '金融街', '望京街'];
const communities = ['阳光小区', '和平花园', '翠湖园', '金色家园', '龙湖花园', '碧桂园'];

const generateRandomAddress = () => {
  const city = cities[Math.floor(Math.random() * cities.length)];
  const district = districts[Math.floor(Math.random() * districts.length)];
  const street = streets[Math.floor(Math.random() * streets.length)];
  const community = communities[Math.floor(Math.random() * communities.length)];
  const building = Math.floor(Math.random() * 20 + 1);
  const unit = Math.floor(Math.random() * 6 + 1);
  const room = Math.floor(Math.random() * 2000 + 101);
  
  formState.address = `${city}${district}${street}${community}${building}号楼${unit}单元${room}室`;
};

// 随机生成面积
const generateRandomArea = () => {
  // 生成 50-300 之间的随机面积，保留2位小数
  formState.area = Number((Math.random() * (300 - 50) + 50).toFixed(2));
};

// 随机生成所有者姓名
const lastNames = ['张', '王', '李', '赵', '刘', '陈', '杨', '黄', '周', '吴'];
const firstNames = ['伟', '芳', '娜', '秀英', '敏', '静', '丽', '强', '磊', '洋', '艳', '勇', '军', '杰', '娟', '涛', '明', '超', '秀兰', '霞'];

const generateRandomOwner = () => {
  const lastName = lastNames[Math.floor(Math.random() * lastNames.length)];
  const firstName = firstNames[Math.floor(Math.random() * firstNames.length)];
  formState.owner = lastName + firstName;
};

// 添加状态筛选的响应式变量
const statusFilter = ref('');

// 监听状态筛选变化
watch(statusFilter, () => {
  // 重置列表和书签
  realEstateList.value = [];
  bookmark.value = '';
  // 重新加载数据
  loadRealEstateList();
});

// 添加复制函数
const handleCopy = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text);
    message.success('已复制到剪贴板');
  } catch (err) {
    message.error('复制失败');
  }
};

// 添加搜索相关的变量和方法
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
    // 当搜索框清空时，重新加载列表
    realEstateList.value = [];
    bookmark.value = '';
    loadRealEstateList();
  }
};

// 初始加载
onMounted(() => {
  loadRealEstateList();
});
</script>

<style scoped>
.realty-agency {
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

.table-footer {
  margin-top: 16px;
  display: flex;
  justify-content: center;
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