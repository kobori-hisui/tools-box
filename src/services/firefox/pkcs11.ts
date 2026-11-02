import { has } from '@/platform/capabilities';

/** browser.pkcs11 - PKCS #11 modules. Firefox only: always guard with isAvailable(). */
export const NAMESPACE = 'pkcs11';

export const isAvailable = (): boolean => has(NAMESPACE);
