import { defineToolMeta } from '@/tools/_core/define-tool';

export default defineToolMeta({
  id: 'json-formatter',
  name: 'JSON Formatter',
  description: 'Format, validate and minify JSON.',
  category: 'developer',
  keywords: ['json', 'format', 'validate'],
  icon: 'wrench',
  version: '0.1.0',
  surfaces: ['popup', 'sidepanel', 'devtools'],
});
