import { defineWxtModule } from 'wxt/modules';

/** Imprime los permisos finales del manifest en cada build (revisa antes de publicar). */
export default defineWxtModule({
  setup(wxt) {
    wxt.hook('build:manifestGenerated', (_, manifest) => {
      wxt.logger.info('[permissions]', manifest.permissions ?? []);
      wxt.logger.info('[host_permissions]', manifest.host_permissions ?? []);
    });
  },
});
