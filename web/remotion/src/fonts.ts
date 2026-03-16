import { loadFont } from '@remotion/google-fonts/Inter';

// Fire-and-forget: @remotion/google-fonts handles delayRender/continueRender internally.
const { fontFamily } = loadFont('normal', {
  weights: ['400', '500', '600', '700', '800'],
  subsets: ['latin'],
});

export const INTER_FONT = fontFamily;
