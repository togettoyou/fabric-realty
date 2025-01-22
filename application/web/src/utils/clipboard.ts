import { message } from 'ant-design-vue';

/**
 * 使用 document.execCommand 进行复制（兼容模式）
 * @param text 要复制的文本
 * @returns 是否复制成功
 */
const fallbackCopyToClipboard = (text: string): boolean => {
  const textArea = document.createElement('textarea');
  textArea.value = text;
  
  // 防止滚动
  textArea.style.top = '0';
  textArea.style.left = '0';
  textArea.style.position = 'fixed';
  textArea.style.opacity = '0';

  document.body.appendChild(textArea);
  textArea.focus();
  textArea.select();

  try {
    const successful = document.execCommand('copy');
    document.body.removeChild(textArea);
    return successful;
  } catch (err) {
    document.body.removeChild(textArea);
    return false;
  }
};

/**
 * 复制文本到剪贴板
 * 优先使用现代的 Clipboard API
 * 如果不支持或者不在安全上下文中，则回退到 execCommand 方法
 */
export const copyToClipboard = async (text: string) => {
  try {
    // 检查是否在安全上下文中并且支持 Clipboard API
    if (window.isSecureContext && navigator.clipboard) {
      await navigator.clipboard.writeText(text);
      message.success('复制成功');
    } else {
      // 回退到兼容模式
      const result = fallbackCopyToClipboard(text);
      if (result) {
        message.success('复制成功');
      } else {
        throw new Error('复制失败');
      }
    }
  } catch (err) {
    console.error('复制失败:', err);
    message.error('复制失败');
  }
};
