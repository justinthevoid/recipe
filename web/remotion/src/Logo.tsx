import { AbsoluteFill } from 'remotion';
import { z } from 'zod';
import { INTER_FONT } from './fonts';
import { BRAND_GRADIENT } from './theme';

export const logoSchema = z.object({
  text: z.string().default('Recipe'),
});

export type LogoProps = z.infer<typeof logoSchema>;

export const Logo: React.FC<LogoProps> = ({ text }) => {
  return (
    <AbsoluteFill
      style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        backgroundColor: 'transparent',
      }}
    >
      <div
        style={{
          fontFamily: INTER_FONT,
          fontSize: 280,
          fontWeight: 800,
          letterSpacing: '-0.03em',
          background: BRAND_GRADIENT,
          backgroundClip: 'text',
          WebkitBackgroundClip: 'text',
          color: 'transparent',
          lineHeight: 1,
        }}
      >
        {text}
      </div>
    </AbsoluteFill>
  );
};
