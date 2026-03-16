// Recipe Aurora Design Tokens — subset of web/src/styles/app.css @theme block

export const AURORA = {
  deepBlue: '#0A2463',
  softBlue: '#8D9FFF',
  cyan: '#3E92CC',
  violet: '#AA6EEE',
  purple: '#7B2D8E',
  softPink: '#F5B9EA',
} as const;

export const SURFACE = {
  canvasBase: '#1A1917',
  elevated: '#2A2825',
  overlay: '#3A3835',
} as const;

export const TEXT = {
  primary: '#FFFFFF',
  secondary: '#ABABAF',
} as const;

export const GLASS = {
  bg: 'rgba(30, 28, 25, 0.55)',
  blur: 20,
  saturate: 1.8,
  border: 'rgba(255, 255, 255, 0.06)',
  borderWidth: 0.5,
  shadow: '0 8px 32px rgba(0, 0, 0, 0.25)',
  insetHighlight: 'inset 0 1px 0 rgba(255, 255, 255, 0.05)',
} as const;

export const BRAND_GRADIENT = 'linear-gradient(135deg, #AA6EEE, #3E92CC)';
export const BRAND_GRADIENT_REVERSE = 'linear-gradient(135deg, #3E92CC, #AA6EEE)';
