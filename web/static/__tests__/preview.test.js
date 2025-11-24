/**
 * Unit tests for preview.js
 * Story: 11-1-css-filter-mapping
 * Epic: 11 (Image Preview System)
 *
 * Tests CSS filter mapping function for 100% coverage (AC-7)
 */

import { recipeToCSSFilters, clamp, isCSSFilterSupported } from '../preview.js';

describe('recipeToCSSFilters', () => {
    // AC-7: Test: Zero parameters → 'none'
    test('returns "none" for zero parameters', () => {
        const recipe = { exposure: 0, contrast: 0, saturation: 0, hue: 0 };
        expect(recipeToCSSFilters(recipe)).toBe('none');
    });

    // AC-7: Test: Exposure +0.5 → 'brightness(150%)'
    test('maps exposure to brightness', () => {
        const recipe = { exposure: 0.5 };
        expect(recipeToCSSFilters(recipe)).toBe('brightness(150%)');
    });

    test('maps negative exposure to brightness', () => {
        const recipe = { exposure: -0.5 };
        expect(recipeToCSSFilters(recipe)).toBe('brightness(50%)');
    });

    // AC-7: Test: Contrast +0.3 → 'contrast(130%)'
    test('maps contrast to contrast', () => {
        const recipe = { contrast: 0.3 };
        expect(recipeToCSSFilters(recipe)).toBe('contrast(130%)');
    });

    test('maps negative contrast to contrast', () => {
        const recipe = { contrast: -0.3 };
        expect(recipeToCSSFilters(recipe)).toBe('contrast(70%)');
    });

    // AC-7: Test: Saturation -0.5 → 'saturate(50%)'
    test('maps saturation to saturate', () => {
        const recipe = { saturation: -0.5 };
        expect(recipeToCSSFilters(recipe)).toBe('saturate(50%)');
    });

    test('maps positive saturation to saturate', () => {
        const recipe = { saturation: 0.5 };
        expect(recipeToCSSFilters(recipe)).toBe('saturate(150%)');
    });

    // AC-7: Test: Hue +30 → 'hue-rotate(30deg)'
    test('maps hue to hue-rotate', () => {
        const recipe = { hue: 30 };
        expect(recipeToCSSFilters(recipe)).toBe('hue-rotate(30deg)');
    });

    test('maps negative hue to hue-rotate', () => {
        const recipe = { hue: -30 };
        expect(recipeToCSSFilters(recipe)).toBe('hue-rotate(-30deg)');
    });

    // AC-7: Test: Temperature +20 → 'sepia(6) hue-rotate(10deg)'
    test('maps temperature to sepia + hue-rotate', () => {
        const recipe = { temperature: 20 };
        expect(recipeToCSSFilters(recipe)).toBe('sepia(6) hue-rotate(10deg)');
    });

    test('maps negative temperature to sepia + hue-rotate', () => {
        const recipe = { temperature: -20 };
        expect(recipeToCSSFilters(recipe)).toBe('sepia(6) hue-rotate(-10deg)');
    });

    // AC-7: Test: Multiple parameters combine correctly
    test('combines multiple filters', () => {
        const recipe = { exposure: 0.5, contrast: 0.3, saturation: -0.5 };
        expect(recipeToCSSFilters(recipe)).toBe('brightness(150%) contrast(130%) saturate(50%)');
    });

    test('combines all supported filters', () => {
        const recipe = {
            exposure: 0.5,
            contrast: 0.3,
            saturation: -0.5,
            hue: 30,
            temperature: 20,
        };
        expect(recipeToCSSFilters(recipe)).toBe('brightness(150%) contrast(130%) saturate(50%) hue-rotate(30deg) sepia(6) hue-rotate(10deg)');
    });

    // AC-7: Test: Out-of-range exposure (+5.0) clamps to 'brightness(200%)'
    test('clamps exposure to 200% maximum', () => {
        const recipe = { exposure: 5.0 };
        expect(recipeToCSSFilters(recipe)).toBe('brightness(200%)');
    });

    // AC-7: Test: Negative exposure (-3.0) clamps to 'brightness(0%)'
    test('clamps exposure to 0% minimum', () => {
        const recipe = { exposure: -3.0 };
        expect(recipeToCSSFilters(recipe)).toBe('brightness(0%)');
    });

    test('clamps contrast to 200% maximum', () => {
        const recipe = { contrast: 5.0 };
        expect(recipeToCSSFilters(recipe)).toBe('contrast(200%)');
    });

    test('clamps contrast to 0% minimum', () => {
        const recipe = { contrast: -3.0 };
        expect(recipeToCSSFilters(recipe)).toBe('contrast(0%)');
    });

    test('clamps saturation to 200% maximum', () => {
        const recipe = { saturation: 5.0 };
        expect(recipeToCSSFilters(recipe)).toBe('saturate(200%)');
    });

    test('clamps saturation to 0% minimum', () => {
        const recipe = { saturation: -3.0 };
        expect(recipeToCSSFilters(recipe)).toBe('saturate(0%)');
    });

    test('clamps hue to 180deg maximum', () => {
        const recipe = { hue: 500 };
        expect(recipeToCSSFilters(recipe)).toBe('hue-rotate(180deg)');
    });

    test('clamps hue to -180deg minimum', () => {
        const recipe = { hue: -500 };
        expect(recipeToCSSFilters(recipe)).toBe('hue-rotate(-180deg)');
    });

    test('clamps temperature to 100 maximum', () => {
        const recipe = { temperature: 500 };
        expect(recipeToCSSFilters(recipe)).toBe('sepia(30) hue-rotate(50deg)');
    });

    test('clamps temperature to -100 minimum', () => {
        const recipe = { temperature: -500 };
        expect(recipeToCSSFilters(recipe)).toBe('sepia(30) hue-rotate(-50deg)');
    });

    // AC-7: Test: Null recipe → 'none'
    test('handles null recipe gracefully', () => {
        expect(recipeToCSSFilters(null)).toBe('none');
    });

    // AC-7: Test: Undefined recipe → 'none'
    test('handles undefined recipe gracefully', () => {
        expect(recipeToCSSFilters(undefined)).toBe('none');
    });

    // AC-7: Test: Invalid parameter types skip gracefully
    test('skips invalid parameters', () => {
        const recipe = { exposure: 'invalid', contrast: 0.3 };
        expect(recipeToCSSFilters(recipe)).toBe('contrast(130%)');
    });

    test('skips null parameters', () => {
        const recipe = { exposure: null, contrast: 0.3, saturation: null };
        expect(recipeToCSSFilters(recipe)).toBe('contrast(130%)');
    });

    test('handles empty recipe object', () => {
        const recipe = {};
        expect(recipeToCSSFilters(recipe)).toBe('none');
    });

    test('handles recipe with only invalid parameters', () => {
        const recipe = { exposure: 'invalid', contrast: null, saturation: undefined };
        expect(recipeToCSSFilters(recipe)).toBe('none');
    });

    // Edge case: Very small values (close to zero but not exactly zero)
    test('applies filter for very small non-zero values', () => {
        const recipe = { exposure: 0.001 };
        expect(recipeToCSSFilters(recipe)).toContain('brightness(');
    });

    // Edge case: Exactly at boundaries
    test('handles exposure exactly at +2.0', () => {
        const recipe = { exposure: 2.0 };
        expect(recipeToCSSFilters(recipe)).toBe('brightness(200%)');
    });

    test('handles exposure exactly at -2.0', () => {
        const recipe = { exposure: -2.0 };
        expect(recipeToCSSFilters(recipe)).toBe('brightness(0%)');
    });

    test('handles hue exactly at +180', () => {
        const recipe = { hue: 180 };
        expect(recipeToCSSFilters(recipe)).toBe('hue-rotate(180deg)');
    });

    test('handles hue exactly at -180', () => {
        const recipe = { hue: -180 };
        expect(recipeToCSSFilters(recipe)).toBe('hue-rotate(-180deg)');
    });
});

describe('clamp', () => {
    test('returns value when within range', () => {
        expect(clamp(50, 0, 100)).toBe(50);
    });

    test('returns min when value below minimum', () => {
        expect(clamp(-10, 0, 100)).toBe(0);
    });

    test('returns max when value above maximum', () => {
        expect(clamp(150, 0, 100)).toBe(100);
    });

    test('handles negative ranges', () => {
        expect(clamp(-5, -10, 10)).toBe(-5);
        expect(clamp(-15, -10, 10)).toBe(-10);
        expect(clamp(15, -10, 10)).toBe(10);
    });

    test('handles floating point values', () => {
        expect(clamp(0.5, 0, 1)).toBe(0.5);
        expect(clamp(-0.5, 0, 1)).toBe(0);
        expect(clamp(1.5, 0, 1)).toBe(1);
    });

    test('handles edge case where min equals max', () => {
        expect(clamp(50, 100, 100)).toBe(100);
    });
});

describe('isCSSFilterSupported', () => {
    test('returns true when CSS.supports exists and filter is supported', () => {
        // Mock CSS.supports
        global.CSS = {
            supports: jest.fn(() => true),
        };

        expect(isCSSFilterSupported()).toBe(true);
        expect(global.CSS.supports).toHaveBeenCalledWith('filter', 'brightness(100%)');
    });

    test('returns false when CSS.supports does not exist', () => {
        global.CSS = undefined;
        expect(isCSSFilterSupported()).toBe(false);
    });

    test('returns false when filter is not supported', () => {
        global.CSS = {
            supports: jest.fn(() => false),
        };

        expect(isCSSFilterSupported()).toBe(false);
    });

    test('returns false when CSS.supports is not a function', () => {
        global.CSS = {
            supports: null,
        };

        expect(isCSSFilterSupported()).toBe(false);
    });
});

// Import additional functions for DOM testing
import { applyPreviewFilter, showDisclaimerHelp, checkBrowserCompatibility } from '../preview.js';

describe('applyPreviewFilter', () => {
    let mockImage;

    beforeEach(() => {
        // Setup DOM
        mockImage = document.createElement('img');
        mockImage.id = 'preview-image';
        document.body.appendChild(mockImage);

        // Mock console.log
        jest.spyOn(console, 'log').mockImplementation(() => {});
        jest.spyOn(console, 'warn').mockImplementation(() => {});
    });

    afterEach(() => {
        // Cleanup
        document.body.innerHTML = '';
        jest.restoreAllMocks();
    });

    test('applies CSS filter to preview image', () => {
        const recipe = { exposure: 0.5, contrast: 0.3 };
        applyPreviewFilter(recipe);

        expect(mockImage.style.filter).toBe('brightness(150%) contrast(130%)');
        expect(console.log).toHaveBeenCalledWith('Preview filter applied: brightness(150%) contrast(130%)');
    });

    test('handles missing preview image element gracefully', () => {
        document.body.innerHTML = ''; // Remove the image
        const recipe = { exposure: 0.5 };

        applyPreviewFilter(recipe);

        expect(console.warn).toHaveBeenCalledWith('Preview image element not found (id="preview-image")');
    });

    test('applies "none" filter for empty recipe', () => {
        const recipe = {};
        applyPreviewFilter(recipe);

        expect(mockImage.style.filter).toBe('none');
    });
});

describe('showDisclaimerHelp', () => {
    beforeEach(() => {
        // Mock window.alert
        global.alert = jest.fn();
    });

    afterEach(() => {
        jest.restoreAllMocks();
    });

    test('shows disclaimer help text in alert', () => {
        showDisclaimerHelp();

        expect(global.alert).toHaveBeenCalledWith(expect.stringContaining('This preview uses CSS filters'));
        expect(global.alert).toHaveBeenCalledWith(expect.stringContaining('Limitations:'));
        expect(global.alert).toHaveBeenCalledWith(expect.stringContaining('Temperature/tint is simplified'));
    });
});

describe('checkBrowserCompatibility', () => {
    beforeEach(() => {
        // Setup DOM
        const modal = document.createElement('div');
        modal.id = 'preview-modal';
        document.body.appendChild(modal);

        // Mock console.warn and alert
        jest.spyOn(console, 'warn').mockImplementation(() => {});
        global.alert = jest.fn();
    });

    afterEach(() => {
        document.body.innerHTML = '';
        jest.restoreAllMocks();
    });

    test('returns true when CSS filters are supported', () => {
        global.CSS = {
            supports: jest.fn(() => true),
        };

        expect(checkBrowserCompatibility()).toBe(true);
        expect(console.warn).not.toHaveBeenCalled();
        expect(global.alert).not.toHaveBeenCalled();
    });

    test('returns false and shows alert when CSS filters not supported', () => {
        global.CSS = {
            supports: jest.fn(() => false),
        };

        expect(checkBrowserCompatibility()).toBe(false);
        expect(console.warn).toHaveBeenCalledWith('CSS filters not supported in this browser');
        expect(global.alert).toHaveBeenCalledWith(expect.stringContaining('Preview not available'));
    });

    test('hides preview modal when CSS filters not supported', () => {
        global.CSS = {
            supports: jest.fn(() => false),
        };

        const modal = document.getElementById('preview-modal');
        checkBrowserCompatibility();

        expect(modal.style.display).toBe('none');
    });
});
