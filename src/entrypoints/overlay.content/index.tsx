import { createRoot } from 'react-dom/client';
import '@/assets/styles/content-ui.css';
import Overlay from './Overlay';

/** Command palette global dentro de un ShadowRoot (estilos aislados). */
export default defineContentScript({
  matches: ['<all_urls>'],
  cssInjectionMode: 'ui',
  async main(ctx) {
    const ui = await createShadowRootUi(ctx, {
      name: 'toolbox-overlay',
      position: 'overlay',
      anchor: 'body',
      onMount: (container) => {
        const app = document.createElement('div');
        container.append(app);
        const root = createRoot(app);
        root.render(<Overlay />);
        return root;
      },
      onRemove: (root) => root?.unmount(),
    });
    ui.mount();
  },
});
