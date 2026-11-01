import { defineToolMeta } from '@/tools/_core/define-tool';

export default defineToolMeta({
  id: 'page-media',
  name: 'Page Media Downloader',
  description: 'List and download images/videos from the current page.',
  category: 'downloaders',
  keywords: ['download', 'images', 'video'],
  icon: 'wrench',
  version: '0.1.0',
  surfaces: ['popup', 'content', 'background'],
  matches: ['<all_urls>'],
});
