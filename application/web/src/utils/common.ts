// 生成UUID
export const generateUUID = () => {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    const r = Math.random() * 16 | 0;
    const v = c === 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
};

// 格式化状态文本
export const getStatusText = (status: string) => {
  switch (status) {
    case 'PENDING':
      return '待完成';
    case 'COMPLETED':
      return '已完成';
    case 'NORMAL':
      return '正常';
    case 'IN_TRANSACTION':
      return '交易中';
    default:
      return '未知';
  }
};

// 获取状态对应的颜色
export const getStatusColor = (status: string) => {
  switch (status) {
    case 'PENDING':
      return 'blue';
    case 'COMPLETED':
    case 'NORMAL':
      return 'green';
    case 'IN_TRANSACTION':
      return 'blue';
    default:
      return 'default';
  }
};

// 格式化金额显示
export const formatPrice = (price: number) => {
  return `¥ ${price}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',');
}; 