# @recipe/ui — Shared Component Library

Svelte 5 components shared between the [web interface](../../web/) and [VSCode extension webview](../../webview/). Provides the core parameter editing UI so both surfaces stay consistent.

## Components

| Component | Description |
|-----------|-------------|
| `ParameterSliderUnit` | Slider control for a single preset parameter with label, value display, and reset |
| `ColorGrading` | Shadow/Midtone/Highlight color wheel controls |
| `ColorBlender` | Color adjustment interface |
| `ToneCurveVisual` | Interactive tone curve display |
| `PhotoPreview` | Image preview with applied adjustments |
| `CollapsibleSection` | Expandable/collapsible panel wrapper |

## Types

Shared TypeScript types for the preset data model:

- `UniversalRecipe` — The intermediate representation used across all conversions
- `ParameterDefinition` — Describes a single editable parameter (range, step, label)
- `ToneCurvePoint` — Point on a tone curve (x, y)
- `ColorAdjustment` — HSL adjustment values
- `ColorGradingZone` — Color grading values for a single tonal zone

## Utilities

- `getNested(obj, path)` / `setNested(obj, path, value)` — Dot-path accessors for nested recipe objects

## Usage

Both consuming packages reference this as a workspace dependency:

```json
{
  "dependencies": {
    "@recipe/ui": "workspace:*"
  }
}
```

Import components directly:

```svelte
<script>
  import { ParameterSliderUnit, ColorGrading } from '@recipe/ui';
</script>
```

## Peer Dependencies

Requires Svelte 5 and Tailwind CSS 4 in the consuming project.
