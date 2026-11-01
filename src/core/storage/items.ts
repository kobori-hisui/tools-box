import { storage } from '#imports';
import { defaultSettings, type Settings } from '@/core/settings/defaults';

export const settingsItem = storage.defineItem<Settings>('sync:settings', {
  fallback: defaultSettings,
  version: 1,
});

export const enabledToolsItem = storage.defineItem<string[]>('local:enabledTools', {
  fallback: [],
  version: 1,
});

export const recentToolsItem = storage.defineItem<string[]>('local:recentTools', {
  fallback: [],
  version: 1,
});

export const favoriteToolsItem = storage.defineItem<string[]>('sync:favoriteTools', {
  fallback: [],
  version: 1,
});
