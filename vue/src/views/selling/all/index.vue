<template>
  <div class="container">
    <el-alert
      type="success"
    >
      <p>账户ID: {{ accountId }}</p>
      <p>用户名: {{ userName }}</p>
      <p>余额: ￥{{ balance }} 元</p>
    </el-alert>
    <div v-if="sellingList.length==0" style="text-align: center;">
      <el-alert
        title="查询不到数据"
        type="warning"
      />
    </div>
    <el-row v-loading="loading" :gutter="20">
      <el-col v-for="(val,index) in sellingList" :key="index" :span="6" :offset="1">
        <el-card class="all-card">
          <div slot="header" class="clearfix">
            <span>{{ val.sellingStatus }}</span>
            <el-button v-if="roles[0] !== 'admin'&&(val.seller===accountId||val.buyer===accountId)&&val.sellingStatus!=='完成'&&val.sellingStatus!=='已过期'&&val.sellingStatus!=='已取消'" style="float: right; padding: 3px 0" type="text" @click="updateSelling(val,'cancelled')">取消</el-button>
            <el-button v-if="roles[0] !== 'admin'&&val.seller===accountId&&val.sellingStatus==='交付中'" style="float: right; padding: 3px 8px" type="text" @click="updateSelling(val,'done')">确认收款</el-button>
            <el-button v-if="roles[0] !== 'admin'&&val.sellingStatus==='销售中'&&val.seller!==accountId" style="float: right; padding: 3px 0" type="text" @click="createSellingByBuy(val)">购买</el-button>
          </div>
          <div class="item">
            <el-tag>房产ID: </el-tag>
            <span>{{ val.objectOfSale }}</span>
          </div>
          <div class="item">
            <el-tag type="success">销售者ID: </el-tag>
            <span>{{ val.seller }}</span>
          </div>
          <div class="item">
            <el-tag type="danger">价格: </el-tag>
            <span>￥{{ val.price }} 元</span>
          </div>
          <div class="item">
            <el-tag type="warning">有效期: </el-tag>
            <span>{{ val.salePeriod }} 天</span>
          </div>
          <div class="item">
            <el-tag type="info">创建时间: </el-tag>
            <span>{{ val.createTime }}</span>
          </div>
          <div class="item">
            <el-tag>购买者ID: </el-tag>
            <span v-if="val.buyer===''">虚位以待</span>
            <span>{{ val.buyer }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { querySellingList, createSellingByBuy, updateSelling } from '@/api/selling'

export default {
  name: 'AllSelling',
  data() {
    return {
      loading: true,
      sellingList: []
    }
  },
  computed: {
    ...mapGetters([
      'accountId',
      'roles',
      'userName',
      'balance'
    ])
  },
  created() {
    querySellingList().then(response => {
      if (response !== null) {
        this.sellingList = response
      }
      this.loading = false
    }).catch(_ => {
      this.loading = false
    })
  },
  methods: {
    createSellingByBuy(item) {
      this.$confirm('是否立即购买?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'success'
      }).then(() => {
        this.loading = true
        createSellingByBuy({
          buyer: this.accountId,
          objectOfSale: item.objectOfSale,
          seller: item.seller
        }).then(response => {
          this.loading = false
          if (response !== null) {
            this.$message({
              type: 'success',
              message: '购买成功!'
            })
          } else {
            this.$message({
              type: 'error',
              message: '购买失败!'
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
          message: '已取消购买'
        })
      })
    },
    updateSelling(item, type) {
      let tip = ''
      if (type === 'done') {
        tip = '确认收款'
      } else {
        tip = '取消操作'
      }
      this.$confirm('是否要' + tip + '?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'success'
      }).then(() => {
        this.loading = true
        updateSelling({
          buyer: item.buyer,
          objectOfSale: item.objectOfSale,
          seller: item.seller,
          status: type
        }).then(response => {
          this.loading = false
          if (response !== null) {
            this.$message({
              type: 'success',
              message: tip + '操作成功!'
            })
          } else {
            this.$message({
              type: 'error',
              message: tip + '操作失败!'
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
          message: '已取消' + tip
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

  .all-card {
    width: 280px;
    height: 380px;
    margin: 18px;
  }
</style>
