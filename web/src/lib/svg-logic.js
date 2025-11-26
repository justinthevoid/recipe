/**
 * SVG Filter Logic
 * Calculates matrices and transfer functions for Lightroom parameter approximation.
 */

import { clamp } from './preview-logic';

// Helper: Multiply two 5x5 matrices (flattened as 20-element arrays, last row implicit 00001)
// A * B means apply B first, then A.
function multiplyMatrices(m1, m2) {
    const result = new Array(20).fill(0);

    // Helper to get element at row r, col c from 5x5 matrix
    // The input arrays are 4x5 (20 elements). The 5th row is implicitly [0, 0, 0, 0, 1]
    const get = (m, r, c) => {
        if (r === 4) return c === 4 ? 1 : 0;
        return m[r * 5 + c];
    };

    for (let r = 0; r < 4; r++) {
        for (let c = 0; c < 5; c++) {
            let sum = 0;
            for (let k = 0; k < 5; k++) {
                sum += get(m1, r, k) * get(m2, k, c);
            }
            result[r * 5 + c] = sum;
        }
    }
    return result;
}

/**
 * Calculate Color Matrix for Temp, Tint, Saturation
 * Returns a 20-element array for feColorMatrix type="matrix"
 */
export function calculateColorMatrix(temp, tint, saturation) {
    // 1. Temperature & Tint Matrix
    // We start with Identity
    const tempTintMatrix = [
        1, 0, 0, 0, 0,
        0, 1, 0, 0, 0,
        0, 0, 1, 0, 0,
        0, 0, 0, 1, 0
    ];

    // Temperature (Blue-Yellow Axis)
    if (temp !== 0) {
        const t = clamp(temp, -100, 100) / 100; // -1.0 to 1.0

        // R channel: Warm (+) increases, Cool (-) decreases
        // B channel: Warm (+) decreases, Cool (-) increases
        // G channel: Slight compensation

        if (t > 0) { // Warm
            tempTintMatrix[0] += t * 0.15;  // R
            tempTintMatrix[6] += t * 0.05;  // G
            tempTintMatrix[12] -= t * 0.15; // B
        } else { // Cool
            tempTintMatrix[0] += t * 0.10;  // R (t is neg)
            tempTintMatrix[12] -= t * 0.20; // B (t is neg, so +=)
        }
    }

    // Tint (Green-Magenta Axis)
    if (tint !== 0) {
        const t = clamp(tint, -100, 100) / 100;
        // Magenta (+): +R, +B, -G
        // Green (-): +G, -R, -B

        tempTintMatrix[6] -= t * 0.15; // G
        tempTintMatrix[0] += t * 0.10; // R
        tempTintMatrix[12] += t * 0.10; // B
    }

    // 2. Saturation Matrix
    let satMatrix = [
        1, 0, 0, 0, 0,
        0, 1, 0, 0, 0,
        0, 0, 1, 0, 0,
        0, 0, 0, 1, 0
    ];

    if (saturation !== 0) {
        const s = 1 + (clamp(saturation, -100, 100) / 100); // 0.0 to 2.0

        // Luminance coefficients (Rec. 709)
        const lumR = 0.2126;
        const lumG = 0.7152;
        const lumB = 0.0722;

        const oneMinusS = 1 - s;

        satMatrix = [
            (lumR * oneMinusS) + s, lumG * oneMinusS, lumB * oneMinusS, 0, 0,
            lumR * oneMinusS, (lumG * oneMinusS) + s, lumB * oneMinusS, 0, 0,
            lumR * oneMinusS, lumG * oneMinusS, (lumB * oneMinusS) + s, 0, 0,
            0, 0, 0, 1, 0
        ];
    }

    // Combine: Apply Saturation FIRST, then Temp/Tint
    // Result = TempTint * Saturation
    const finalMatrix = multiplyMatrices(tempTintMatrix, satMatrix);

    return finalMatrix.join(' ');
}

// sRGB <-> Linear conversions for physically accurate exposure
function sRGBtoLinear(x) {
    return x <= 0.04045 ? x / 12.92 : Math.pow((x + 0.055) / 1.055, 2.4);
}

function linearToSRGB(x) {
    return x <= 0.0031308 ? x * 12.92 : 1.055 * Math.pow(x, 1.0 / 2.4) - 0.055;
}


/**
 * Calculate Component Transfer Table for Contrast/Exposure/Tone/Curve
 * Returns tableValues string for feFuncR/G/B
 */
export function calculateTransferTable(exposure, contrast, highlights, shadows, whites, blacks, curveHighlights, curveLights, curveDarks, curveShadows) {
    const steps = 50;
    const values = [];

    const exp = exposure || 0;
    const cont = (contrast || 0) / 100; // -1.0 to 1.0

    // Tone params (-100 to 100) -> (-1.0 to 1.0)
    const high = (highlights || 0) / 100;
    const shad = (shadows || 0) / 100;
    const wht = (whites || 0) / 100;
    const blk = (blacks || 0) / 100;

    // Curve params (-100 to 100) -> (-1.0 to 1.0)
    const cHigh = (curveHighlights || 0) / 100;
    const cLight = (curveLights || 0) / 100;
    const cDark = (curveDarks || 0) / 100;
    const cShad = (curveShadows || 0) / 100;

    const exposureMult = Math.pow(2, exp);
    const slopeFactor = Math.pow(2, cont);
    const k = 4 * slopeFactor;

    for (let i = 0; i <= steps; i++) {
        let x = i / steps; // 0.0 to 1.0 (sRGB signal)

        // 1. Convert to Linear Space
        let lin = sRGBtoLinear(x);

        // 2. Apply Exposure
        lin = lin * exposureMult;

        // 3. Convert back to sRGB
        let res = linearToSRGB(lin);

        // 4. Apply Contrast (Logistic Sigmoid)
        if (cont !== 0) {
            const f = (val) => 1 / (1 + Math.exp(-k * (val - 0.5)));
            const min = f(0);
            const max = f(1);
            res = (f(res) - min) / (max - min);
        }

        // 5. Apply Basic Tone Mapping (Approximate)
        // We apply a curve that affects specific ranges

        // Shadows/Blacks (affect lower end)
        if (shad !== 0 || blk !== 0) {
            // Simple quadratic/cubic influence for shadows
            // Range 0.0 to 0.5 mostly
            const shadowFactor = (1 - res) * (1 - res);
            res += shad * 0.2 * shadowFactor * Math.sin(res * Math.PI);

            // Blacks affect the very bottom
            const blackFactor = Math.pow(1 - res, 4);
            res += blk * 0.15 * blackFactor;
        }

        // Highlights/Whites (affect upper end)
        if (high !== 0 || wht !== 0) {
            // Range 0.5 to 1.0 mostly
            const highlightFactor = res * res;
            res += high * 0.2 * highlightFactor * Math.sin(res * Math.PI);

            // Whites affect the very top
            const whiteFactor = Math.pow(res, 4);
            res += wht * 0.15 * whiteFactor;
        }

        // 6. Apply Parametric Curve (Approximate)
        if (cShad !== 0) {
            // Shadows: 0.0 - 0.25
            const factor = Math.max(0, 1 - res * 4); // Linear falloff from 0 to 0.25
            res += cShad * 0.15 * factor;
        }
        if (cDark !== 0) {
            // Darks: 0.25 - 0.5 (Peak at 0.375)
            const factor = Math.max(0, 1 - Math.abs(res - 0.375) * 4);
            res += cDark * 0.15 * factor;
        }
        if (cLight !== 0) {
            // Lights: 0.5 - 0.75 (Peak at 0.625)
            const factor = Math.max(0, 1 - Math.abs(res - 0.625) * 4);
            res += cLight * 0.15 * factor;
        }
        if (cHigh !== 0) {
            // Highlights: 0.75 - 1.0
            const factor = Math.max(0, (res - 0.75) * 4);
            res += cHigh * 0.15 * factor;
        }

        // Clamp
        values.push(Math.max(0, Math.min(1, res)));
    }

    return values.join(' ');
}
