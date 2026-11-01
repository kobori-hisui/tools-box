import { loadToolImpl, getToolsForUrl } from '@/tools/_core/registry';

/**
 * Content script UNICO que carga bajo demanda el content.ts de cada
 * herramienta cuyo `matches` coincide con la URL (soporta SPAs).
 */
export default defineContentScript({
  matches: ['<all_urls>'],
  runAt: 'document_idle',
  async main(ctx) {
    const mounted = new Set<string>();

    const mountFor = async (url: string) => {
      for (const meta of getToolsForUrl(url)) {
        if (mounted.has(meta.id)) continue;
        mounted.add(meta.id);
        const impl = await loadToolImpl(meta.id);
        await impl?.onContent?.({ ctx, meta, url });
      }
    };

    await mountFor(location.href);
    ctx.addEventListener(window, 'wxt:locationchange', ({ newUrl }) => {
      void mountFor(newUrl.href);
    });
  },
});
