<template>
  <div class="container">
    <el-alert
      type="success"
    >
      <p>账户ID: {{ accountId }}</p>
      <p>用户名: {{ userName }}</p>
      <p>余额: ￥{{ balance }} 元</p>
      <p>当发起出售、捐赠或质押操作后，担保状态为true</p>
      <p>当担保状态为false时，才可发起出售、捐赠或质押操作</p>
    </el-alert>
    <div v-if="realEstateList.length==0" style="text-align: center;">
      <el-alert
        title="查询不到数据"
        type="warning"
      />
    </div>
    <el-row v-loading="loading" :gutter="20">
      <el-col v-for="(val,index) in realEstateList" :key="index" :span="6" :offset="1">
        <el-card class="realEstate-card">
          <div slot="header" class="clearfix">
            担保状态:
            <span style="color: rgb(255, 0, 0);">{{ val.encumbrance }}</span>
          </div>

          <div class="item">
            <el-tag>房产ID: </el-tag>
            <span>{{ val.realEstateId }}</span>
          </div>
          <div class="item">
            <el-tag type="success">业主ID: </el-tag>
            <span>{{ val.proprietor }}</span>
          </div>
          <div class="item">
            <el-tag type="warning">总空间: </el-tag>
            <span>{{ val.totalArea }} ㎡</span>
          </div>
          <div class="item">
            <el-tag type="danger">居住空间: </el-tag>
            <span>{{ val.livingSpace }} ㎡</span>
          </div>

          <div v-if="!val.encumbrance&&roles[0] !== 'admin'">
            <el-button type="text" @click="openDialog(val)">出售</el-button>
            <el-divider direction="vertical" />
            <el-button type="text" @click="openDonatingDialog(val)">捐赠</el-button>
          </div>
          <el-rate v-if="roles[0] === 'admin'" />
        </el-card>
      </el-col>
    </el-row>
    <el-dialog v-loading="loadingDialog" :visible.sync="dialogCreateSelling" :close-on-click-modal="false" @close="resetForm('realForm')">
      <el-form ref="realForm" :model="realForm" :rules="rules" label-width="100px">
        <el-form-item label="价格 (元)" prop="price">
          <el-input-number v-model="realForm.price" :precision="2" :step="10000" :min="0" />
        </el-form-item>
        <el-form-item label="有效期 (天)" prop="salePeriod">
          <el-input-number v-model="realForm.salePeriod" :min="1" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button type="primary" @click="createSelling('realForm')">立即出售</el-button>
        <el-button @click="dialogCreateSelling = false">取 消</el-button>
      </div>
    </el-dialog>
    <el-dialog v-loading="loadingDialog" :visible.sync="dialogCreateDonating" :close-on-click-modal="false" @close="resetForm('DonatingForm')">
      <el-form ref="DonatingForm" :model="DonatingForm" :rules="rulesDonating" label-width="100px">
        <el-form-item label="业主" prop="proprietor">
          <el-select v-model="DonatingForm.proprietor" placeholder="请选择业主" @change="selectGet">
            <el-option
              v-for="item in accountList"
              :key="item.accountId"
              :label="item.userName"
              :value="item.accountId"
            >
              <span style="float: left">{{ item.userName }}</span>
              <span style="float: right; color: #8492a6; font-size: 13px">{{ item.accountId }}</span>
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button type="primary" @click="createDonating('DonatingForm')">立即捐赠</el-button>
        <el-button @click="dialogCreateDonating = false">取 消</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { queryAccountList } from '@/api/account'
import { queryRealEstateList } from '@/api/realEstate'
import { createSelling } from '@/api/selling'
import { createDonating } from '@/api/donating'

export default {
  name: 'RealeState',
  data() {
    var checkArea = (rule, value, callback) => {
      if (value <= 0) {
        callback(new Error('必须大于0'))
      } else {
        callback()
      }
    }
    return {
      loading: true,
      loadingDialog: false,
      realEstateList: [],
      dialogCreateSelling: false,
      dialogCreateDonating: false,
      realForm: {
        price: 0,
        salePeriod: 0
      },
      rules: {
        price: [
          { validator: checkArea, trigger: 'blur' }
        ],
        salePeriod: [
          { validator: checkArea, trigger: 'blur' }
        ]
      },
      DonatingForm: {
        proprietor: ''
      },
      rulesDonating: {
        proprietor: [
          { required: true, message: '请选择业主', trigger: 'change' }
        ]
      },
      accountList: [],
      valItem: {}
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
    if (this.roles[0] === 'admin') {
      queryRealEstateList().then(response => {
        if (response !== null) {
          this.realEstateList = response
        }
        this.loading = false
      }).catch(_ => {
        this.loading = false
      })
    } else {
      queryRealEstateList({ proprietor: this.accountId }).then(response => {
        if (response !== null) {
          this.realEstateList = response
        }
        this.loading = false
      }).catch(_ => {
        this.loading = false
      })
    }
  },
  methods: {
    openDialog(item) {
      this.dialogCreateSelling = true
      this.valItem = item
    },
    openDonatingDialog(item) {
      this.dialogCreateDonating = true
      this.valItem = item
      queryAccountList().then(response => {
        if (response !== null) {
          // 过滤掉管理员和当前用户
          this.accountList = response.filter(item =>
            item.userName !== '管理员' && item.accountId !== this.accountId
          )
        }
      })
    },
    createSelling(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.$confirm('是否立即出售?', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'success'
          }).then(() => {
            this.loadingDialog = true
            createSelling({
              objectOfSale: this.valItem.realEstateId,
              seller: this.valItem.proprietor,
              price: this.realForm.price,
              salePeriod: this.realForm.salePeriod
            }).then(response => {
              this.loadingDialog = false
              this.dialogCreateSelling = false
              if (response !== null) {
                this.$message({
                  type: 'success',
                  message: '出售成功!'
                })
              } else {
                this.$message({
                  type: 'error',
                  message: '出售失败!'
                })
              }
              setTimeout(() => {
                window.location.reload()
              }, 1000)
            }).catch(_ => {
              this.loadingDialog = false
              this.dialogCreateSelling = false
            })
          }).catch(() => {
            this.loadingDialog = false
            this.dialogCreateSelling = false
            this.$message({
              type: 'info',
              message: '已取消出售'
            })
          })
        } else {
          return false
        }
      })
    },
    createDonating(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.$confirm('是否立即捐赠?', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'success'
          }).then(() => {
            this.loadingDialog = true
            createDonating({
              objectOfDonating: this.valItem.realEstateId,
              donor: this.valItem.proprietor,
              grantee: this.DonatingForm.proprietor
            }).then(response => {
              this.loadingDialog = false
              this.dialogCreateDonating = false
              if (response !== null) {
                this.$message({
                  type: 'success',
                  message: '捐赠成功!'
                })
              } else {
                this.$message({
                  type: 'error',
                  message: '捐赠失败!'
                })
              }
              setTimeout(() => {
                window.location.reload()
              }, 1000)
            }).catch(_ => {
              this.loadingDialog = false
              this.dialogCreateDonating = false
            })
          }).catch(() => {
            this.loadingDialog = false
            this.dialogCreateDonating = false
            this.$message({
              type: 'info',
              message: '已取消捐赠'
            })
          })
        } else {
          return false
        }
      })
    },
    resetForm(formName) {
      this.$refs[formName].resetFields()
    },
    selectGet(accountId) {
      this.DonatingForm.proprietor = accountId
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

  .realEstate-card {
    width: 280px;
    height: 340px;
    margin: 18px;
  }
</style>
