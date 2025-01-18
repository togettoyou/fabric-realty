.content {
  margin: 24px auto;
  max-width: 1200px;
  padding: 0 24px;
}

// ... existing code ...

const columns = [
  {
    title: '房产ID',
    dataIndex: 'id',
    key: 'id',
    width: 180,
    ellipsis: true,
    customCell: () => ({
      style: { cursor: 'copy' },
      onClick: (e: MouseEvent) => {
        const text = (e.target as HTMLElement).innerText;
        navigator.clipboard.writeText(text);
        message.success('已复制到剪贴板');
      },
    }),
  },
  {
    title: '地址',
    dataIndex: 'address',
    key: 'address',
    width: 200,
    ellipsis: true,
  },
  {
    title: '面积',
    dataIndex: 'area',
    key: 'area',
    width: 120,
    align: 'right' as const,
  },
  {
    title: '所有者',
    dataIndex: 'owner',
    key: 'owner',
    width: 120,
    ellipsis: true,
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    width: 100,
  },
  {
    title: '创建时间',
    dataIndex: 'createTime',
    key: 'createTime',
    width: 180,
  },
  {
    title: '更新时间',
    dataIndex: 'updateTime',
    key: 'updateTime',
    width: 180,
  },
]; 