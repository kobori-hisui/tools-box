export interface Settings {
  theme: 'system' | 'light' | 'dark';
  language: string;
  telemetry: boolean;
}

export const defaultSettings: Settings = {
  theme: 'system',
  language: 'en',
  telemetry: false,
};
