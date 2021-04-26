import request from '@/utils/request'

// 查询捐赠列表(可查询所有，也可根据发起捐赠人查询)
export function queryDonatingList(data) {
  return request({
    url: '/queryDonatingList',
    method: 'post',
    data
  })
}

// 根据受赠人(受赠人AccountId)查询捐赠(受赠的)(供受赠人查询)
export function queryDonatingListByGrantee(data) {
  return request({
    url: '/queryDonatingListByGrantee',
    method: 'post',
    data
  })
}

// 更新捐赠状态（确认受赠、取消） Status取值为 完成"done"、取消"cancelled"
export function updateDonating(data) {
  return request({
    url: '/updateDonating',
    method: 'post',
    data
  })
}

// 发起捐赠
export function createDonating(data) {
  return request({
    url: '/createDonating',
    method: 'post',
    data
  })
}
