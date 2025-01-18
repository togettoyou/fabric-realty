import axios from 'axios';
import type { ApiResponse } from '../types';

const request = axios.create({
  baseURL: '/api',
  timeout: 10000,
});

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    const res = response.data as ApiResponse<any>;
    if (res.code === 200) {
      return res.data;
    }
    return Promise.reject(new Error(res.message || '请求失败'));
  },
  (error) => {
    // 处理 HTTP 错误
    if (error.response) {
      const res = error.response.data as ApiResponse<any>;
      return Promise.reject(new Error(res.message || '服务器错误'));
    }
    // 处理请求超时
    if (error.code === 'ECONNABORTED') {
      return Promise.reject(new Error('请求超时，请稍后重试'));
    }
    // 处理网络错误
    if (!window.navigator.onLine) {
      return Promise.reject(new Error('网络连接已断开，请检查网络'));
    }
    // 其他错误
    return Promise.reject(new Error(error.message || '请求失败'));
  }
);

export default request; 