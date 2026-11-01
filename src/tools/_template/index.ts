import { defineTool } from '@/tools/_core/define-tool';
import Panel from './ui/Panel';

export default defineTool({
  Panel,
  onBackground: (ctx) => import('./background').then((m) => m.default(ctx)),
  onContent: (ctx) => import('./content').then((m) => m.default(ctx)),
});
