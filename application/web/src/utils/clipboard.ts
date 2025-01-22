import { message } from 'ant-design-vue';

export const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text);
    message.success('复制成功');
  } catch (err) {
    message.error('复制失败');
  }
};
