import axios from 'axios'
import {
  MessageBox,
  Message
} from 'element-ui'

const service = axios.create({
  baseURL: process.env.VUE_APP_BASE_API,
  timeout: 5000
})

service.interceptors.response.use(
  response => {
    const res = response.data
    if (res.code !== 200) {
      MessageBox.alert('服务器开小差了', 'error', {
        confirmButtonText: '确定',
        type: 'warning'
      })
      return Promise.reject(new Error(res.msg || 'Error'))
    } else {
      return res.data
    }
  },
  error => {
    console.log('err' + error)
    Message({
      message: error.message,
      type: 'error',
      duration: 5 * 1000
    })
    return Promise.reject(error)
  }
)

export default service
