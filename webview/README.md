# Recipe — VSCode Extension Webview

The Svelte 5 UI that runs inside the VSCode extension's webview panel. Provides the parameter editor, photo preview, and all user interaction for the [extension](../extension/).

## How It Fits Together

```
VSCode Extension Host (TypeScript)
  ↕ postMessage IPC
Webview (this package — Svelte 5)
  ↕ WASM calls
recipe.wasm (Go — photo preview filters)
```

The webview communicates with the extension host via `postMessage` for file operations (open, save, patch parameters). Photo preview filters run locally in the webview via the WASM module.

## Build Constraints

The webview output is a **single JS file and single CSS file** — no code splitting. This is required by VSCode's Content Security Policy, which restricts webview script sources. Vite is configured with `build.rollupOptions` to produce a single chunk.

Output location: `extension/dist/webview/webview.{js,css}`

## Development

```bash
cd webview
bun install
bun run dev      # Watch mode (rebuilds on change)
bun run build    # Production build
bun run check    # Svelte type checking
bun run test     # Run tests (Vitest)
```

During development, run the webview watcher alongside the extension watcher, then press F5 in VSCode to launch the Extension Development Host:

```bash
# Terminal 1
cd webview && bun run dev

# Terminal 2
cd extension && bun run dev

# Then F5 in VSCode
```

## Structure

```
webview/
├── src/
│   ├── App.svelte                 # Root component (editor + preview modes)
│   ├── lib/
│   │   ├── components/            # Parameter editors, dropdowns, preview
│   │   ├── components/ui/         # Base UI primitives (bits-ui based)
│   │   ├── state/np3.svelte.ts    # Svelte 5 runes state management
│   │   ├── vscode.ts              # VSCode API bindings (postMessage)
│   │   └── wasm.svelte.ts         # WASM module loader
│   └── app.css                    # Tailwind CSS entry point
├── public/
│   └── recipe.wasm                # WASM binary (built separately)
└── vite.config.ts                 # Single-file output config
```

## Shared Components

Core editing components (sliders, color grading, tone curves) come from `@recipe/ui` in [`packages/ui/`](../packages/ui/). This package adds VSCode-specific wiring: IPC communication, dark mode integration, toast notifications, and the two-mode layout (Editor vs. Preview).

## Testing

Tests use Vitest with `@testing-library/svelte` and jsdom. Run with `bun run test`.
