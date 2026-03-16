import { AbsoluteFill, useCurrentFrame, useVideoConfig } from 'remotion';
import { NoiseTexture } from './NoiseTexture';
import { AURORA, SURFACE } from './theme';

interface BlobConfig {
  color: string;
  width: number;
  height: number;
  blur: number;
  x: number; // 0-1 fraction of canvas width
  y: number; // 0-1 fraction of canvas height
  opacity: number;
  moveX: number;
  moveY: number;
}

const BLOBS: BlobConfig[] = [
  { color: AURORA.deepBlue, width: 500, height: 500, blur: 120, x: 0.15, y: 0.3, opacity: 0.7, moveX: 60, moveY: 50 },
  { color: AURORA.violet, width: 400, height: 400, blur: 100, x: 0.45, y: 0.6, opacity: 0.5, moveX: -60, moveY: -50 },
  { color: AURORA.cyan, width: 450, height: 450, blur: 110, x: 0.75, y: 0.35, opacity: 0.45, moveX: 60, moveY: -50 },
  { color: AURORA.purple, width: 350, height: 350, blur: 90, x: 0.85, y: 0.7, opacity: 0.4, moveX: -60, moveY: 50 },
  { color: AURORA.softBlue, width: 380, height: 380, blur: 100, x: 0.3, y: 0.8, opacity: 0.35, moveX: 60, moveY: -50 },
  { color: AURORA.softPink, width: 300, height: 300, blur: 80, x: 0.6, y: 0.15, opacity: 0.3, moveX: -60, moveY: 50 },
];

export const AuroraBackground: React.FC<{ children?: React.ReactNode }> = ({ children }) => {
  const frame = useCurrentFrame();
  const { fps, width, height } = useVideoConfig();

  const breathingDuration = 7 * fps;
  const breathingProgress = (frame % breathingDuration) / breathingDuration;
  const breathingValue = Math.sin(breathingProgress * Math.PI * 2);

  return (
    <AbsoluteFill style={{ backgroundColor: SURFACE.canvasBase }}>
      {BLOBS.map((blob, index) => {
        const offsetX = breathingValue * blob.moveX;
        const offsetY = breathingValue * blob.moveY;

        return (
          <div
            key={index}
            style={{
              position: 'absolute',
              width: blob.width,
              height: blob.height,
              borderRadius: '50%',
              background: blob.color,
              filter: `blur(${blob.blur}px)`,
              opacity: blob.opacity,
              left: blob.x * width - blob.width / 2 + offsetX,
              top: blob.y * height - blob.height / 2 + offsetY,
            }}
          />
        );
      })}

      {/* Warm-gray overlay for depth */}
      <div
        style={{
          position: 'absolute',
          inset: 0,
          background: 'rgba(26, 25, 23, 0.12)',
        }}
      />

      <NoiseTexture opacity={0.03} baseFrequency={0.4} />

      {children}
    </AbsoluteFill>
  );
};
