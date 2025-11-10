// format-selector.js - Target format selection
// Epic 2, Story 2-5: Target Format Selection
// Allows user to choose output format (NP3, XMP, lrtemplate, Costyle, Costylepack) with smart defaults

let sourceFormat = null;
let targetFormat = null;

/**
 * Format definitions with metadata
 */
const FORMATS = {
    np3: {
        name: 'Nikon Picture Control',
        extension: '.np3',
        software: 'Nikon cameras (D850, Z9, etc.)',
        description: 'Native preset format for Nikon cameras. Load directly in camera settings.',
        badgeClass: 'badge-blue',
    },
    xmp: {
        name: 'Lightroom CC Preset',
        extension: '.xmp',
        software: 'Adobe Lightroom CC, Lightroom Mobile',
        description: 'Modern Lightroom preset format. Works with cloud sync.',
        badgeClass: 'badge-purple',
    },
    lrtemplate: {
        name: 'Lightroom Classic Preset',
        extension: '.lrtemplate',
        software: 'Adobe Lightroom Classic (desktop)',
        description: 'Legacy Lightroom preset format. For Lightroom Classic 7.3 and earlier.',
        badgeClass: 'badge-teal',
    },
    costyle: {
        name: 'Capture One Style',
        extension: '.costyle',
        software: 'Capture One Pro',
        description: 'Capture One style format. Import as adjustments or styles.',
        badgeClass: 'badge-orange',
    },
    costylepack: {
        name: 'Capture One Style Bundle',
        extension: '.costylepack',
        software: 'Capture One Pro',
        description: 'ZIP bundle of Capture One styles. Import multiple styles at once.',
        badgeClass: 'badge-orange',
    },
};

/**
 * Smart default target format based on source format
 */
const SMART_DEFAULTS = {
    np3: 'xmp',           // Nikon users want Lightroom CC
    xmp: 'np3',           // Lightroom users want Nikon
    lrtemplate: 'xmp',    // Lightroom Classic users want Lightroom CC
    costyle: 'xmp',       // Capture One users want Lightroom CC
    costylepack: 'xmp',   // Capture One bundle users want Lightroom CC
};

/**
 * Display format selection UI
 * @param {string} detectedFormat - Source format detected in Story 2-3
 */
export function displayFormatSelector(detectedFormat) {
    if (!detectedFormat || !FORMATS[detectedFormat]) {
        throw new Error('Invalid source format');
    }

    sourceFormat = detectedFormat;
    targetFormat = SMART_DEFAULTS[detectedFormat];

    renderFormatSelector();
}

/**
 * Render format selection UI
 */
function renderFormatSelector() {
    const container = document.getElementById('conversionControls');
    if (!container) {
        console.error('Conversion controls container not found');
        return;
    }

    let html = `
        <div class="format-selector">
            <h3>Convert to:</h3>
            <div class="format-options">
    `;

    for (const [formatKey, formatData] of Object.entries(FORMATS)) {
        const isDisabled = formatKey === sourceFormat;
        const isSelected = formatKey === targetFormat;
        const disabledClass = isDisabled ? 'disabled' : '';
        const selectedClass = isSelected ? 'selected' : '';

        html += `
            <div class="format-option ${disabledClass} ${selectedClass}"
                 data-format="${formatKey}"
                 ${isDisabled ? 'title="Cannot convert to same format"' : ''}>
                <input type="radio"
                       id="format-${formatKey}"
                       name="targetFormat"
                       value="${formatKey}"
                       ${isSelected ? 'checked' : ''}
                       ${isDisabled ? 'disabled' : ''}
                       class="format-radio">
                <label for="format-${formatKey}" class="format-label">
                    <div class="format-header">
                        <span class="format-badge ${formatData.badgeClass}">
                            ${formatData.name}
                        </span>
                        <span class="format-extension">${formatData.extension}</span>
                    </div>
                    <div class="format-software">${formatData.software}</div>
                    <div class="format-description">${formatData.description}</div>
                </label>
            </div>
        `;
    }

    html += `
            </div>
            <button id="convertButton" class="convert-button">
                Convert to ${FORMATS[targetFormat].name}
            </button>
            <button id="downloadButton" class="convert-button" style="display: none;" disabled>
                Download
            </button>
            <div id="downloadStatus" class="status" style="display: none;"></div>
            <div id="downloadError" class="error-message" style="display: none;"></div>
        </div>
    `;

    container.innerHTML = html;
    container.style.display = 'block';

    // Add event listeners
    attachFormatSelectorListeners();

    console.log(`Format selector displayed: ${sourceFormat} → ${targetFormat} (default)`);
}

/**
 * Attach event listeners to format options
 */
function attachFormatSelectorListeners() {
    // Radio button change events
    const radioButtons = document.querySelectorAll('.format-radio');
    radioButtons.forEach(radio => {
        radio.addEventListener('change', handleFormatChange);
    });

    // Format option click events (click anywhere on the option)
    const formatOptions = document.querySelectorAll('.format-option:not(.disabled)');
    formatOptions.forEach(option => {
        option.addEventListener('click', () => {
            const format = option.dataset.format;
            if (format !== sourceFormat) {
                targetFormat = format;
                renderFormatSelector();
            }
        });
    });

    // Convert button event (handled in Story 2-6)
    const convertButton = document.getElementById('convertButton');
    if (convertButton) {
        convertButton.addEventListener('click', handleConvertClick);
    }
}

/**
 * Handle format selection change
 */
function handleFormatChange(event) {
    const newFormat = event.target.value;
    if (newFormat !== sourceFormat) {
        targetFormat = newFormat;
        updateConvertButton();
        dispatchFormatSelectedEvent(newFormat);
    }
}

/**
 * Update convert button text
 */
function updateConvertButton() {
    const convertButton = document.getElementById('convertButton');
    if (convertButton) {
        convertButton.textContent = `Convert to ${FORMATS[targetFormat].name}`;
    }
}

/**
 * Handle convert button click (Story 2-6 will implement actual conversion)
 */
function handleConvertClick() {
    console.log(`Convert button clicked: ${sourceFormat} → ${targetFormat}`);

    // Dispatch event for Story 2-6
    dispatchConvertRequestEvent(sourceFormat, targetFormat);
}

/**
 * Dispatch format selected event
 */
function dispatchFormatSelectedEvent(format) {
    const event = new CustomEvent('formatSelected', {
        detail: { format }
    });
    window.dispatchEvent(event);
    console.log(`formatSelected event dispatched: ${format}`);
}

/**
 * Dispatch convert request event (for Story 2-6)
 */
function dispatchConvertRequestEvent(fromFormat, toFormat) {
    const event = new CustomEvent('convertRequest', {
        detail: { fromFormat, toFormat }
    });
    window.dispatchEvent(event);
    console.log(`convertRequest event dispatched: ${fromFormat} → ${toFormat}`);
}

/**
 * Get selected target format
 */
export function getTargetFormat() {
    return targetFormat;
}

/**
 * Get source format
 */
export function getSourceFormat() {
    return sourceFormat;
}

/**
 * Clear format selector
 */
export function clearFormatSelector() {
    const container = document.getElementById('conversionControls');
    if (container) {
        container.innerHTML = '';
        container.style.display = 'none';
    }
    sourceFormat = null;
    targetFormat = null;
    console.log('Format selector cleared');
}
