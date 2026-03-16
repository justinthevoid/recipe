import { AbsoluteFill } from 'remotion';
import { z } from 'zod';
import { AuroraBackground } from './AuroraBackground';
import { INTER_FONT } from './fonts';
import { AURORA, BRAND_GRADIENT, GLASS, TEXT } from './theme';

export const ogImageSchema = z.object({
  headline: z.string().default('Nikon \u2194 Lightroom'),
  subtitle: z.string().default('Free browser-based preset converter. Private. Instant.'),
  sourceFormat: z.string().default('NP3'),
  targetFormat: z.string().default('XMP'),
  siteUrl: z.string().default('recipe.shuttercoach.app'),
});

export type OGImageProps = z.infer<typeof ogImageSchema>;

export const OGImage: React.FC<OGImageProps> = ({
  headline,
  subtitle,
  sourceFormat,
  targetFormat,
  siteUrl,
}) => {
  return (
    <AbsoluteFill>
      <AuroraBackground />

      {/* Content overlay */}
      <AbsoluteFill
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          padding: 60,
        }}
      >
        {/* Glass card */}
        <div
          style={{
            background: 'rgba(30, 28, 25, 0.65)',
            border: `${GLASS.borderWidth}px solid ${GLASS.border}`,
            borderRadius: 24,
            boxShadow: `${GLASS.shadow}, ${GLASS.insetHighlight}`,
            padding: 60,
            display: 'flex',
            flexDirection: 'column',
            gap: 24,
            width: '100%',
            maxWidth: 1000,
            position: 'relative',
          }}
        >
          {/* Recipe logo text — top left */}
          <div
            style={{
              fontFamily: INTER_FONT,
              fontSize: 32,
              fontWeight: 800,
              letterSpacing: '-0.03em',
              background: BRAND_GRADIENT,
              backgroundClip: 'text',
              WebkitBackgroundClip: 'text',
              color: 'transparent',
              lineHeight: 1,
            }}
          >
            Recipe
          </div>

          {/* Headline */}
          <div
            style={{
              fontFamily: INTER_FONT,
              fontSize: 52,
              fontWeight: 800,
              color: TEXT.primary,
              letterSpacing: -1,
              lineHeight: 1.15,
            }}
          >
            {headline}
          </div>

          {/* Format badge row */}
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: 16,
            }}
          >
            <div
              style={{
                background: AURORA.violet,
                borderRadius: 8,
                padding: '6px 16px',
                fontFamily: INTER_FONT,
                fontSize: 18,
                fontWeight: 600,
                color: TEXT.primary,
              }}
            >
              {sourceFormat}
            </div>
            <div
              style={{
                fontFamily: INTER_FONT,
                fontSize: 24,
                color: TEXT.primary,
              }}
            >
              {'\u2192'}
            </div>
            <div
              style={{
                background: AURORA.cyan,
                borderRadius: 8,
                padding: '6px 16px',
                fontFamily: INTER_FONT,
                fontSize: 18,
                fontWeight: 600,
                color: TEXT.primary,
              }}
            >
              {targetFormat}
            </div>
          </div>

          {/* Subtitle */}
          <div
            style={{
              fontFamily: INTER_FONT,
              fontSize: 18,
              fontWeight: 400,
              color: TEXT.secondary,
              lineHeight: 1.4,
            }}
          >
            {subtitle}
          </div>

          {/* Site URL — bottom right */}
          <div
            style={{
              fontFamily: INTER_FONT,
              fontSize: 14,
              fontWeight: 500,
              color: TEXT.secondary,
              opacity: 0.6,
              alignSelf: 'flex-end',
              marginTop: 8,
            }}
          >
            {siteUrl}
          </div>
        </div>
      </AbsoluteFill>
    </AbsoluteFill>
  );
};
