# Recipe — Web Interface

The public-facing website and interactive preset converter at [recipe.shuttercoach.app](https://recipe.shuttercoach.app). Built with Astro and Svelte 5, all processing runs client-side via WebAssembly.

## Tech Stack

- **Framework:** Astro 6 (static site generation)
- **Interactive components:** Svelte 5 (embedded in Astro pages)
- **Styling:** Tailwind CSS 4
- **State management:** nanostores
- **Conversion engine:** Go compiled to WebAssembly

## Structure

```
web/
├── src/
│   ├── pages/index.astro      # Landing page
│   ├── layouts/Layout.astro   # Base HTML layout with SEO meta
│   ├── components/            # Astro + Svelte components
│   │   ├── ConversionCard.svelte   # Drag-and-drop file conversion
│   │   ├── EditorView.svelte       # Full parameter editor
│   │   ├── AuroraBackground.svelte # Animated background
│   │   ├── Explainer.astro         # How-it-works section
│   │   ├── FAQ.astro               # Frequently asked questions
│   │   └── ...
│   └── lib/
│       ├── wasm.ts            # WASM module initialization
│       ├── converter.ts       # Conversion function wrappers
│       ├── format-detector.ts # NP3/XMP file detection
│       └── stores/            # Shared state (nanostores)
├── public/
│   ├── recipe.wasm            # Go WASM binary (built in CI)
│   └── images/                # Demo and OG images
├── astro.config.mjs           # Astro config with Svelte + sitemap
└── wrangler.jsonc             # Cloudflare Pages deployment config
```

## Development

```bash
cd web
bun install
bun run dev        # Start dev server with hot reload
bun run build      # Production build to dist/
bun run preview    # Preview production build locally
```

The WASM binary must be built separately before the converter will work locally:

```bash
# From repository root
make wasm
```

## Deployment

Deployed automatically to Cloudflare Pages on push to `main` via GitHub Actions. The workflow builds the WASM binary, runs `astro build`, and deploys `web/dist/`.

## Shared Components

Interactive UI components (sliders, color grading, tone curves) are imported from `@recipe/ui` — the shared component library in [`packages/ui/`](../packages/ui/).
