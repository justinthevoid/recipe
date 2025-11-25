/**
 * SVG Filter Logic
 * Calculates matrices and transfer functions for Lightroom parameter approximation.
 */

import { clamp } from './preview-logic';

/**
 * Calculate Color Matrix for Temp, Tint, Saturation
 * Returns a 20-element array for feColorMatrix type="matrix"
 */
export function calculateColorMatrix(temp, tint, saturation) {
    // Identity matrix
    let r = [1, 0, 0, 0, 0];
    let g = [0, 1, 0, 0, 0];
    let b = [0, 0, 1, 0, 0];
    let a = [0, 0, 0, 1, 0];

    // 1. Temperature (Blue-Yellow Axis)
    // Range: -100 to +100
    // Warm (+): +Red, -Blue
    // Cool (-): +Blue, -Red
    if (temp !== 0) {
        const t = clamp(temp, -100, 100) / 100; // -1.0 to 1.0

        // Warm
        if (t > 0) {
            r[0] += t * 0.2;  // Boost Red
            b[2] -= t * 0.2;  // Reduce Blue
            g[1] += t * 0.05; // Slight Green boost for warmth
        }
        // Cool
        else {
            r[0] += t * 0.1;  // Reduce Red (t is negative)
            b[2] -= t * 0.2;  // Boost Blue (t is negative, so -= is +)
        }
    }

    // 2. Tint (Green-Magenta Axis)
    // Range: -100 to +100
    // Magenta (+): +Red, +Blue, -Green
    // Green (-): +Green, -Red, -Blue
    if (tint !== 0) {
        const t = clamp(tint, -100, 100) / 100; // -1.0 to 1.0

        g[1] -= t * 0.2; // +Tint = Less Green, -Tint = More Green
        r[0] += t * 0.1; // +Tint = More Red
        b[2] += t * 0.1; // +Tint = More Blue
    }

    // 3. Saturation
    // Range: -100 to +100
    // We use a simplified saturation matrix multiplication
    if (saturation !== 0) {
        const s = 1 + (clamp(saturation, -100, 100) / 100); // 0.0 to 2.0

        // Luminance coefficients (Rec. 709)
        const lumR = 0.2126;
        const lumG = 0.7152;
        const lumB = 0.0722;

        const oneMinusS = 1 - s;

        const satMatrix = [
            (lumR * oneMinusS) + s, lumG * oneMinusS, lumB * oneMinusS, 0, 0,
            lumR * oneMinusS, (lumG * oneMinusS) + s, lumB * oneMinusS, 0, 0,
            lumR * oneMinusS, lumG * oneMinusS, (lumB * oneMinusS) + s, 0, 0,
            0, 0, 0, 1, 0
        ];

        // Multiply current matrix by saturation matrix
        // Note: This is a simplified application, ideally we'd do full matrix multiplication
        // But for preview speed, applying saturation to the diagonal is a decent approximation
        // combined with the luminance weights.

        // Let's just return the saturation matrix combined with our temp/tint shifts
        // (Simplified: applying saturation on top of the shifts)
        r[0] = satMatrix[0] + (r[0] - 1);
        g[1] = satMatrix[6] + (g[1] - 1);
        b[2] = satMatrix[12] + (b[2] - 1);
    }

    // Flatten to string
    return [
        r.join(' '),
        g.join(' '),
        b.join(' '),
        a.join(' ')
    ].join(' ');
}

/**
 * Calculate Component Transfer Table for Contrast/Exposure
 * Returns tableValues string for feFuncR/G/B
 */
export function calculateTransferTable(exposure, contrast) {
    // Exposure: -5.0 to +5.0 (Linear multiplier)
    // Contrast: -100 to +100 (S-Curve)

    const steps = 20;
    const values = [];

    // Normalize inputs
    const exp = exposure || 0;
    const cont = (contrast || 0) / 100; // -1.0 to 1.0

    // Pre-calculate exposure multiplier (2^exposure)
    // +1 EV = 2x brightness, -1 EV = 0.5x brightness
    const exposureMult = Math.pow(2, exp);

    for (let i = 0; i <= steps; i++) {
        let x = i / steps; // 0.0 to 1.0

        // 1. Apply Contrast (S-Curve)
        // Simple S-curve function: f(x) = x + c * x * (1 - x) * (x - 0.5) ? No that's cubic
        // Let's use a cosine-based S-curve or simple power function
        // If contrast > 0: Darker shadows, brighter highlights
        if (cont !== 0) {
            // Center around 0.5
            // Formula: x + contrast * (x - 0.5) * x * (1-x) * factor
            // This is a rough approximation.
            // Better: (x - 0.5) * (contrast + 1) + 0.5 for linear contrast
            // For S-curve:
            if (cont > 0) {
                // Increase contrast
                x = (x - 0.5) * (1 + cont) + 0.5;
                // Clamp
                x = Math.max(0, Math.min(1, x));
            } else {
                // Decrease contrast (linear compression towards gray)
                x = (x - 0.5) * (1 + cont) + 0.5;
            }
        }

        // 2. Apply Exposure
        x = x * exposureMult;

        // Clamp final value
        values.push(Math.max(0, Math.min(1, x)));
    }

    return values.join(' ');
}
