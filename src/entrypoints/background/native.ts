import { onMessage } from '@/core/messaging';
import { NativeClient, NativeError } from '@/core/native';
import { supports } from '@/platform/capabilities';

let client: NativeClient | undefined;

/**
 * Bridges extension pages / content scripts to the native host. Pages never
 * talk to the host directly: the service worker owns the single port.
 */
export function registerNative() {
  onMessage('native:call', async ({ data }) => {
    if (!supports('nativeMessaging')) {
      throw new NativeError('unavailable', 'nativeMessaging permission not granted');
    }
    client ??= new NativeClient();
    return client.call(data.method, data.params);
  });

  onMessage('native:available', async () => {
    if (!supports('nativeMessaging')) return false;
    try {
      client ??= new NativeClient({ timeoutMs: 3000 });
      await client.call('ping');
      return true;
    } catch {
      return false;
    }
  });
}
