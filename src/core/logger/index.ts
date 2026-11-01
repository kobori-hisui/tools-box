export const logger = {
  debug: (...a: unknown[]) => import.meta.env.DEV && console.debug('[toolbox]', ...a),
  info: (...a: unknown[]) => console.info('[toolbox]', ...a),
  warn: (...a: unknown[]) => console.warn('[toolbox]', ...a),
  error: (...a: unknown[]) => console.error('[toolbox]', ...a),
};
