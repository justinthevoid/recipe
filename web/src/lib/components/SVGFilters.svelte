<script>
    import { calculateColorMatrix, calculateTransferTable } from "../svg-logic";

    export let parameters = {};

    // Helper to safely get parameter values
    const getVal = (...keys) => {
        if (!parameters) return 0;
        for (const key of keys) {
            if (
                parameters[key] != null &&
                typeof parameters[key] === "number"
            ) {
                return parameters[key];
            }
        }
        return 0;
    };

    // Reactive values
    $: temp = getVal("Temperature");
    $: tint = getVal("Tint");
    $: saturation = getVal("Saturation");
    $: exposure = getVal("Exposure", "Exposure2012");
    $: contrast = getVal("Contrast", "Contrast2012");

    // Calculate filter attributes
    $: colorMatrix = calculateColorMatrix(temp, tint, saturation);
    $: transferTable = calculateTransferTable(exposure, contrast);
</script>

<svg style="position: absolute; width: 0; height: 0; overflow: hidden;">
    <defs>
        <filter id="preview-filter" color-interpolation-filters="sRGB">
            <!-- 1. Color Matrix (Temp, Tint, Saturation) -->
            <feColorMatrix
                type="matrix"
                values={colorMatrix}
                result="colored"
            />

            <!-- 2. Component Transfer (Exposure, Contrast) -->
            <feComponentTransfer in="colored" result="final">
                <feFuncR type="table" tableValues={transferTable} />
                <feFuncG type="table" tableValues={transferTable} />
                <feFuncB type="table" tableValues={transferTable} />
            </feComponentTransfer>
        </filter>
    </defs>
</svg>
