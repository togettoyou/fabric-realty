// API 响应基础结构
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

// 分页查询结果
export interface PageResult<T> {
  bookmark: string;
  fetchedRecordsCount: number;
  records: T[];
  recordsCount: number;
}

// 房产信息
export interface RealEstate {
  id: string;
  propertyAddress: string;
  area: number;
  currentOwner: string;
  status: 'NORMAL' | 'IN_TRANSACTION';
  createTime: string;
  updateTime: string;
}

// 交易信息
export interface Transaction {
  id: string;
  realEstateId: string;
  seller: string;
  buyer: string;
  price: number;
  status: 'PENDING' | 'COMPLETED' | 'CANCELLED';
  createTime: string;
  updateTime: string;
}

// 房产列表查询结果
export type RealEstatePageResult = PageResult<RealEstate>;

// 交易列表查询结果
export type TransactionPageResult = PageResult<Transaction>; 