import { defineWxtModule } from 'wxt/modules';
import { existsSync, readdirSync } from 'node:fs';
import { resolve } from 'node:path';

/** Avisa en build si una herramienta no tiene meta.ts o index.ts. */
export default defineWxtModule({
  setup(wxt) {
    wxt.hook('ready', () => {
      const root = resolve(wxt.config.srcDir, 'tools/categories');
      if (!existsSync(root)) return;
      for (const cat of readdirSync(root, { withFileTypes: true })) {
        if (!cat.isDirectory()) continue;
        for (const tool of readdirSync(resolve(root, cat.name), { withFileTypes: true })) {
          if (!tool.isDirectory()) continue;
          const dir = resolve(root, cat.name, tool.name);
          for (const file of ['meta.ts', 'index.ts']) {
            if (!existsSync(resolve(dir, file))) {
              wxt.logger.warn(`[tools] ${cat.name}/${tool.name} no tiene ${file}`);
            }
          }
        }
      }
    });
  },
});
