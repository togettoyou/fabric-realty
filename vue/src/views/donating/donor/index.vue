<template>
  <div class="container">
    <el-alert
      type="success"
    >
      <p>账户ID: {{ accountId }}</p>
      <p>用户名: {{ userName }}</p>
      <p>余额: ￥{{ balance }} 元</p>
    </el-alert>
    <div v-if="donatingList.length==0" style="text-align: center;">
      <el-alert
        title="查询不到数据"
        type="warning"
      />
    </div>
    <el-row v-loading="loading" :gutter="20">
      <el-col v-for="(val,index) in donatingList" :key="index" :span="6" :offset="1">
        <el-card class="d-me-card">
          <div slot="header" class="clearfix">
            <span>{{ val.donatingStatus }}</span>
            <el-button v-if="val.donatingStatus === '捐赠中'" style="float: right; padding: 3px 0" type="text" @click="updateDonating(val)">取消</el-button>
          </div>
          <div class="item">
            <el-tag>房产ID: </el-tag>
            <span>{{ val.objectOfDonating }}</span>
          </div>
          <div class="item">
            <el-tag type="success">捐赠者ID: </el-tag>
            <span>{{ val.donor }}</span>
          </div>
          <div class="item">
            <el-tag type="danger">受赠人ID: </el-tag>
            <span>{{ val.grantee }}</span>
          </div>
          <div class="item">
            <el-tag type="warning">创建时间: </el-tag>
            <span>{{ val.createTime }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { queryDonatingList, updateDonating } from '@/api/donating'

export default {
  name: 'DonatingDonor',
  data() {
    return {
      loading: true,
      donatingList: []
    }
  },
  computed: {
    ...mapGetters([
      'accountId',
      'userName',
      'balance'
    ])
  },
  created() {
    queryDonatingList({ donor: this.accountId }).then(response => {
      if (response !== null) {
        this.donatingList = response
      }
      this.loading = false
    }).catch(_ => {
      this.loading = false
    })
  },
  methods: {
    updateDonating(item) {
      this.$confirm('是否要取消捐赠?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'success'
      }).then(() => {
        this.loading = true
        updateDonating({
          donor: item.donor,
          grantee: item.grantee,
          objectOfDonating: item.objectOfDonating,
          status: 'cancelled'
        }).then(response => {
          this.loading = false
          if (response !== null) {
            this.$message({
              type: 'success',
              message: '操作成功!'
            })
          } else {
            this.$message({
              type: 'error',
              message: '操作失败!'
            })
          }
          setTimeout(() => {
            window.location.reload()
          }, 1000)
        }).catch(_ => {
          this.loading = false
        })
      }).catch(() => {
        this.loading = false
        this.$message({
          type: 'info',
          message: '已取消操作'
        })
      })
    }
  }
}

</script>

<style>
  .container{
    width: 100%;
    text-align: center;
    min-height: 100%;
    overflow: hidden;
  }
  .tag {
    float: left;
  }

  .item {
    font-size: 14px;
    margin-bottom: 18px;
    color: #999;
  }

  .clearfix:before,
  .clearfix:after {
    display: table;
  }
  .clearfix:after {
    clear: both
  }

  .d-me-card {
    width: 280px;
    height: 300px;
    margin: 18px;
  }
</style>
