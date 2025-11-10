// parameter-display.js - Display preset parameters

let currentParameters = null;
let currentFormat = null;
let isPanelExpanded = true;

/**
 * Extract and display parameters from preset file
 * @param {Uint8Array} fileData - Raw file bytes
 * @param {string} format - Detected format ("np3" | "xmp" | "lrtemplate")
 */
export async function displayParameters(fileData, format) {
    if (!fileData || !format) {
        throw new Error('File data and format required');
    }

    // Check if WASM is ready
    if (typeof extractParameters !== 'function') {
        throw new Error('WASM module not loaded');
    }

    console.log(`Extracting parameters from ${format} file...`);
    const startTime = performance.now();

    try {
        // Call WASM function (returns Promise<string> containing JSON)
        const jsonString = await extractParameters(fileData, format);
        const parameters = JSON.parse(jsonString);

        const elapsedTime = performance.now() - startTime;
        console.log(`Parameters extracted: ${Object.keys(parameters).length} params (${elapsedTime.toFixed(2)}ms)`);

        // Store for later use
        currentParameters = parameters;
        currentFormat = format;

        // Display in UI
        renderParameterPanel(parameters, format);

        return parameters;

    } catch (error) {
        console.error('Parameter extraction failed:', error);
        throw new Error(`Unable to extract parameters: ${error.message || error}`);
    }
}

/**
 * Render parameter panel in UI
 */
function renderParameterPanel(parameters, format) {
    const container = document.getElementById('parameterPanel');
    if (!container) {
        console.error('Parameter panel container not found');
        return;
    }

    // Group parameters by category
    const grouped = groupParameters(parameters, format);

    // Build HTML
    let html = `
        <div class="parameter-panel ${isPanelExpanded ? 'expanded' : 'collapsed'}">
            <div class="parameter-header">
                <h3>Parameters</h3>
                <button id="toggleParameters" class="toggle-button" type="button">
                    ${isPanelExpanded ? 'Hide' : 'Show'}
                </button>
            </div>
    `;

    if (isPanelExpanded) {
        for (const [category, params] of Object.entries(grouped)) {
            html += `
                <div class="parameter-section">
                    <h4>${category}</h4>
                    <div class="parameter-grid">
            `;

            for (const [name, value] of Object.entries(params)) {
                const displayValue = formatParameterValue(value);
                html += `
                    <div class="parameter-row">
                        <span class="parameter-name">${name}</span>
                        <span class="parameter-value">${displayValue}</span>
                    </div>
                `;
            }

            html += `
                    </div>
                </div>
            `;
        }
    }

    html += '</div>';

    container.innerHTML = html;
    container.style.display = 'block';

    // Add event listener for toggle button
    const toggleButton = document.getElementById('toggleParameters');
    if (toggleButton) {
        toggleButton.addEventListener('click', toggleParameterPanel);
    }
}

/**
 * Group parameters by category
 */
function groupParameters(parameters, format) {
    const groups = {
        'Basic Adjustments': {},
        'Color Adjustments': {},
        'Detail Adjustments': {},
    };

    // Basic adjustments
    const basicParams = ['Exposure', 'Exposure2012', 'Contrast', 'Contrast2012',
                        'Highlights', 'Highlights2012', 'Shadows', 'Shadows2012',
                        'Whites', 'Whites2012', 'Blacks', 'Blacks2012'];

    // Color adjustments
    const colorParams = ['Vibrance', 'Saturation', 'Temperature', 'Tint'];

    // Detail adjustments
    const detailParams = ['Clarity', 'Clarity2012', 'Sharpness', 'Dehaze',
                         'Texture', 'GrainAmount', 'GrainSize'];

    for (const [key, value] of Object.entries(parameters)) {
        if (value === null || value === undefined) continue;

        // Skip Name parameter (it's metadata, not a photo edit parameter)
        if (key === 'Name') continue;

        if (basicParams.includes(key)) {
            groups['Basic Adjustments'][key] = value;
        } else if (colorParams.includes(key)) {
            groups['Color Adjustments'][key] = value;
        } else if (detailParams.includes(key)) {
            groups['Detail Adjustments'][key] = value;
        }
    }

    // Remove empty groups
    for (const [category, params] of Object.entries(groups)) {
        if (Object.keys(params).length === 0) {
            delete groups[category];
        }
    }

    return groups;
}

/**
 * Format parameter value for display
 */
function formatParameterValue(value) {
    if (value === null || value === undefined) {
        return '—';
    }

    if (typeof value === 'number') {
        // Format numbers with sign (+ for positive, - for negative)
        if (value > 0) {
            return `+${value.toFixed(2)}`;
        } else if (value < 0) {
            return value.toFixed(2);
        } else {
            return '0';
        }
    }

    return String(value);
}

/**
 * Toggle parameter panel expanded/collapsed
 */
function toggleParameterPanel() {
    isPanelExpanded = !isPanelExpanded;

    // Re-render with current state
    if (currentParameters && currentFormat) {
        renderParameterPanel(currentParameters, currentFormat);
    }
}

/**
 * Clear parameter panel
 */
export function clearParameterPanel() {
    const container = document.getElementById('parameterPanel');
    if (container) {
        container.innerHTML = '';
        container.style.display = 'none';
    }
    currentParameters = null;
    currentFormat = null;
}

/**
 * Get current parameters
 */
export function getCurrentParameters() {
    return currentParameters;
}
