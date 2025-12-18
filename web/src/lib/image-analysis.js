/**
 * Image Analysis Logic
 * Calculates histograms and statistics from image data.
 */

/**
 * Analyze an image element to get histogram data
 * @param {HTMLImageElement} imgElement 
 * @returns {Promise<{histogram: {r: number[], g: number[], b: number[], l: number[]}, stats: any}>}
 */
export async function analyzeImage(imgElement) {
    return new Promise((resolve, reject) => {
        try {
            const canvas = document.createElement('canvas');
            const ctx = canvas.getContext('2d');

            // Limit analysis size for performance (e.g., max 1000px width)
            const MAX_WIDTH = 1000;
            let width = imgElement.naturalWidth;
            let height = imgElement.naturalHeight;

            if (width > MAX_WIDTH) {
                const ratio = MAX_WIDTH / width;
                width = MAX_WIDTH;
                height = height * ratio;
            }

            canvas.width = width;
            canvas.height = height;

            ctx.drawImage(imgElement, 0, 0, width, height);

            const imageData = ctx.getImageData(0, 0, width, height);
            const data = imageData.data;

            // Initialize histograms
            const rHist = new Array(256).fill(0);
            const gHist = new Array(256).fill(0);
            const bHist = new Array(256).fill(0);
            const lHist = new Array(256).fill(0);

            let totalPixels = width * height;
            let maxCount = 0;

            for (let i = 0; i < data.length; i += 4) {
                const r = data[i];
                const g = data[i + 1];
                const b = data[i + 2];
                // Rec. 709 Luminance
                const l = Math.round(0.2126 * r + 0.7152 * g + 0.0722 * b);

                rHist[r]++;
                gHist[g]++;
                bHist[b]++;
                lHist[l]++;
            }

            // Normalize histograms to 0-1 range relative to peak? 
            // Or just return raw counts. Raw counts are better for flexibility.

            resolve({
                histogram: {
                    r: rHist,
                    g: gHist,
                    b: bHist,
                    l: lHist
                },
                width,
                height,
                totalPixels
            });
        } catch (e) {
            reject(e);
        }
    });
}

/**
 * Calculate Auto-Tone Exposure Offset
 * Analyzes luminance histogram to find the exposure shift needed to center the midtones
 * or expand the dynamic range.
 */
export function calculateAutoExposure(lHist, totalPixels) {
    // 1. Find the median luminance
    let count = 0;
    let median = 0;
    const halfPixels = totalPixels / 2;

    for (let i = 0; i < 256; i++) {
        count += lHist[i];
        if (count >= halfPixels) {
            median = i;
            break;
        }
    }

    // Target median is usually around 128 (middle gray)
    // But for "Auto Exposure", we often look at the "key" of the image.
    // A simple approach: Shift median to 128.

    // Calculate EV shift needed.
    // Current brightness = median / 255
    // Target brightness = 0.5
    // Exposure factor = Target / Current
    // EV = log2(factor)

    // Safety check for very dark/bright images
    const safeMedian = Math.max(10, Math.min(245, median));
    const currentVal = safeMedian / 255;
    const targetVal = 0.5;

    const factor = targetVal / currentVal;
    let evShift = Math.log2(factor);

    // Clamp extreme shifts
    evShift = Math.max(-3, Math.min(3, evShift));

    return evShift;
}
