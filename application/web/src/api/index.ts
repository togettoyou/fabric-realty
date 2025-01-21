import request from '../utils/request';
import type { RealEstatePageResult, TransactionPageResult, RealEstate, Transaction } from '../types';

// 区块数据类型
interface BlockData {
  block_num: number;
  block_hash: string;
  data_hash: string;
  prev_hash: string;
  tx_count: number;
  save_time: string;
}

// 区块查询结果类型
interface BlockQueryResult {
  blocks: BlockData[];
  total: number;
  page_size: number;
  page_num: number;
  has_more: boolean;
}

// 不动产登记机构接口
export const realtyAgencyApi = {
  // 创建房产信息
  createRealEstate: (data: {
    id: string;
    address: string;
    area: number;
    owner: string;
  }) => request.post<never, void>('/realty-agency/realty/create', data),

  // 查询房产信息
  getRealEstate: (id: string) => request.get<never, RealEstate>(`/realty-agency/realty/${id}`),

  // 分页查询房产列表
  getRealEstateList: (params: { pageSize: number; bookmark: string; status?: string }) =>
    request.get<never, RealEstatePageResult>('/realty-agency/realty/list', { params }),

  // 分页查询区块列表
  getBlockList: (params: { pageSize?: number; pageNum?: number }) =>
    request.get<never, BlockQueryResult>('/realty-agency/block/list', { params }),
};

// 交易平台接口
export const tradingPlatformApi = {
  // 生成交易
  createTransaction: (data: {
    txId: string;
    realEstateId: string;
    seller: string;
    buyer: string;
    price: number;
  }) => request.post<never, void>('/trading-platform/transaction/create', data),

  // 查询房产信息
  getRealEstate: (id: string) => request.get<never, RealEstate>(`/trading-platform/realty/${id}`),

  // 查询交易信息
  getTransaction: (txId: string) => request.get<never, Transaction>(`/trading-platform/transaction/${txId}`),

  // 分页查询交易列表
  getTransactionList: (params: { pageSize: number; bookmark: string; status?: string }) =>
    request.get<never, TransactionPageResult>('/trading-platform/transaction/list', { params }),

  // 分页查询区块列表
  getBlockList: (params: { pageSize?: number; pageNum?: number }) =>
    request.get<never, BlockQueryResult>('/trading-platform/block/list', { params }),
};

// 银行接口
export const bankApi = {
  // 完成交易
  completeTransaction: (txId: string) =>
    request.post<never, void>(`/bank/transaction/complete/${txId}`),

  // 查询交易信息
  getTransaction: (txId: string) => request.get<never, Transaction>(`/bank/transaction/${txId}`),

  // 分页查询交易列表
  getTransactionList: (params: { pageSize: number; bookmark: string; status?: string }) =>
    request.get<never, TransactionPageResult>('/bank/transaction/list', { params }),

  // 分页查询区块列表
  getBlockList: (params: { pageSize?: number; pageNum?: number }) =>
    request.get<never, BlockQueryResult>('/bank/block/list', { params }),
};
