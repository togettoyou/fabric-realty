import request from '../utils/request';
import type { RealEstatePageResult, TransactionPageResult, RealEstate, Transaction } from '../types';

// 房产相关接口
export const realtyApi = {
  // 创建房产信息
  createRealEstate: (data: {
    id: string;
    address: string;
    area: number;
    owner: string;
  }) => request.post<never, void>('/realty-agency/realty/create', data),

  // 查询房产信息
  getRealEstate: (id: string) => request.get<never, RealEstate>(`/query/realty/${id}`),

  // 分页查询房产列表
  getRealEstateList: (params: { pageSize: number; bookmark: string }) =>
    request.get<never, RealEstatePageResult>('/query/realty/list', { params }),
};

// 交易相关接口
export const transactionApi = {
  // 创建交易
  createTransaction: (data: {
    txId: string;
    realEstateId: string;
    seller: string;
    buyer: string;
    price: number;
  }) => request.post<never, void>('/trading-platform/transaction/create', data),

  // 完成交易
  completeTransaction: (txId: string) =>
    request.post<never, void>(`/bank/transaction/complete/${txId}`),

  // 查询交易信息
  getTransaction: (txId: string) => request.get<never, Transaction>(`/query/transaction/${txId}`),

  // 分页查询交易列表
  getTransactionList: (params: { pageSize: number; bookmark: string }) =>
    request.get<never, TransactionPageResult>('/query/transaction/list', { params }),
}; 