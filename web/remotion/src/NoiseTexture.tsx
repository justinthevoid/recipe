import type React from 'react';
import { useId } from 'react';

export const NoiseTexture: React.FC<{
  opacity?: number;
  baseFrequency?: number;
  octaves?: number;
  blendMode?: 'overlay' | 'soft-light' | 'multiply' | 'normal';
  seed?: number;
}> = ({
  opacity = 0.03,
  baseFrequency = 0.4,
  octaves = 4,
  blendMode = 'overlay',
  seed = 0,
}) => {
  const filterId = `noise-${useId().replace(/:/g, '')}`;

  return (
    <svg
      style={{
        position: 'absolute',
        inset: 0,
        width: '100%',
        height: '100%',
        pointerEvents: 'none',
        opacity,
        mixBlendMode: blendMode,
      }}
    >
      <filter id={filterId}>
        <feTurbulence
          type="fractalNoise"
          baseFrequency={baseFrequency}
          numOctaves={octaves}
          seed={seed}
          stitchTiles="stitch"
        />
      </filter>
      <rect width="100%" height="100%" filter={`url(#${filterId})`} />
    </svg>
  );
};
