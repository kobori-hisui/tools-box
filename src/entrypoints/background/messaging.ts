import { onMessage } from '@/core/messaging';
import { getAllTools } from '@/tools/_core/registry';

export function registerMessaging() {
  onMessage('tools:list', () => getAllTools());
}
