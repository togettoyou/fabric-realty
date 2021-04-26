import request from '@/utils/request'

// 获取登录界面角色选择列表
export function queryAccountList() {
  return request({
    url: '/queryAccountList',
    method: 'post'
  })
}

// 登录
export function login(data) {
  return request({
    url: '/queryAccountList',
    method: 'post',
    data
  })
}
