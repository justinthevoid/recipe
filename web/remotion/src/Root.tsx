import './index.css';
import { Folder, Still } from 'remotion';
import { Logo, logoSchema } from './Logo';
import { OGImage, ogImageSchema } from './OGImage';

export const RemotionRoot = () => (
  <Folder name="Branding">
    <Still
      id="Logo"
      component={Logo}
      width={2400}
      height={600}
      schema={logoSchema}
      defaultProps={{ text: 'Recipe' }}
    />
    <Still
      id="OGImage"
      component={OGImage}
      width={1200}
      height={630}
      schema={ogImageSchema}
      defaultProps={{
        headline: 'Nikon \u2194 Lightroom',
        subtitle: 'Free browser-based preset converter. Private. Instant.',
        sourceFormat: 'NP3',
        targetFormat: 'XMP',
        siteUrl: 'recipe.shuttercoach.app',
      }}
    />
  </Folder>
);
