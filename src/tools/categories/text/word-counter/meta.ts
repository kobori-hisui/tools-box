import { defineToolMeta } from '@/tools/_core/define-tool';

export default defineToolMeta({
  id: 'word-counter',
  name: 'Word Counter',
  description: 'Count words, characters and reading time.',
  category: 'text',
  keywords: ['words', 'count', 'characters'],
  icon: 'wrench',
  version: '0.1.0',
  surfaces: ['popup', 'sidepanel'],
});
