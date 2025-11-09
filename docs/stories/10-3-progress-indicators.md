# Story 10.3: Progress Indicators for Batch Conversions

Status: ready-for-dev

## Story

As a **photographer converting multiple preset files**,
I want **clear visual feedback on conversion progress for each file and overall batch status**,
so that **I can monitor conversion progress, identify errors quickly, and know when downloads are ready**.

## Acceptance Criteria

**AC-1: Batch Progress Indicator**
- ✅ Overall batch progress displayed above file grid
- ✅ Progress text format: "Converting 3 of 10..." (dynamic count)
- ✅ Progress updates in real-time as files complete
- ✅ Progress bar visual indicator (0-100% filled bar)
- ✅ Completion message: "All 10 files converted successfully" (or "8 of 10 converted, 2 failed")
- ✅ Progress indicator hidden when no conversions in progress

**AC-2: Per-File Status Indicators**
- ✅ **Queued state** (initial):
  - Icon: ⏱️ (clock emoji)
  - Text: "Queued"
  - Color: Gray (#999)
  - File card background: Default (white)
- ✅ **Processing state** (during conversion):
  - Icon: ⏳ (animated spinner or hourglass)
  - Text: "Converting..."
  - Color: Blue (#0073E6)
  - File card background: Light blue tint (#F0F7FF)
  - Spinner animation (CSS `@keyframes` rotation)
- ✅ **Complete state** (conversion success):
  - Icon: ✓ (checkmark, green)
  - Text: "Complete"
  - Color: Green (#4CAF50)
  - File card background: Light green tint (#E8F5E9)
  - Download button visible
- ✅ **Error state** (conversion failure):
  - Icon: ✕ (X, red)
  - Text: "Error: <specific error message>"
  - Color: Red (#D32F2F)
  - File card background: Light red tint (#FFEBEE)
  - Error message displayed (e.g., "Invalid file format", "Parsing failed")

**AC-3: Smooth State Transitions**
- ✅ Transitions between states use CSS transitions (0.3s ease)
- ✅ No janky updates or layout shifts during state changes
- ✅ Background color fades smoothly (0.3s transition)
- ✅ Icon changes instantly (no animation delay)
- ✅ Status text updates without flicker

**AC-4: Visual Feedback During Processing**
- ✅ Spinner animation during "processing" state:
  - Rotating hourglass icon (360° rotation, 1s duration, infinite loop)
  - Or CSS loading spinner (border-radius circle with animated border)
- ✅ Progress bar for batch processing (optional enhancement):
  - Horizontal bar showing percentage (0-100%)
  - Filled portion blue (#0073E6), unfilled light gray (#E0E0E0)
- ✅ Processing state visually distinct (color change, animation)

**AC-5: Completion State with Download Buttons**
- ✅ When file conversion completes, download button appears in file card
- ✅ Download button label: "Download <format>" (e.g., "Download XMP")
- ✅ Download button styled as primary action (blue background, white text)
- ✅ Clicking download button triggers file download (browser download)
- ✅ Downloaded filename: `<original-name>_converted.<ext>` (e.g., `preset_converted.xmp`)
- ✅ Download button persists after download (user can download again)

**AC-6: Error State with Specific Error Messages**
- ✅ Error message displayed per file (not global error)
- ✅ Common error messages:
  - "Invalid file format" (file content doesn't match extension)
  - "Parsing failed: <reason>" (NP3/XMP/lrtemplate parse error)
  - "Conversion failed: <reason>" (mapping or generation error)
  - "Unsupported format" (file extension mismatch)
- ✅ Error message truncated if >50 characters (full message on hover)
- ✅ Error state persistent until user removes file or retries
- ✅ Retry button shown for error state (optional): "Retry Conversion"

**AC-7: Cancel In-Progress Batch Conversions**
- ✅ "Cancel All" button visible during batch processing
- ✅ Clicking "Cancel All" stops all pending conversions:
  - Files in "queued" state → revert to queued
  - Files in "processing" state → abort conversion (if possible)
  - Completed files → remain complete (downloads available)
- ✅ Confirmation dialog: "Cancel all in-progress conversions?" (Yes/No)
- ✅ After cancellation, batch progress indicator updates: "Conversion cancelled (3 of 10 complete)"
- ✅ "Cancel All" button hidden when no conversions in progress

## Tasks / Subtasks

### Task 1: Create Batch Progress Indicator Component (AC-1)
- [ ] Add batch progress HTML to `web/index.html` (above file grid):
  ```html
  <div class="batch-progress" id="batch-progress" hidden>
    <div class="batch-progress__header">
      <h3 class="batch-progress__title" id="batch-progress-title">Converting 0 of 0...</h3>
      <button class="batch-progress__cancel button button--secondary" id="cancel-batch">Cancel All</button>
    </div>
    <div class="batch-progress__bar">
      <div class="batch-progress__fill" id="batch-progress-fill" style="width: 0%;"></div>
    </div>
  </div>
  ```
- [ ] Add batch progress styles to `web/css/components.css`:
  ```css
  .batch-progress {
    background: white;
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: 1.5rem;
    margin-bottom: 2rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .batch-progress__header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }

  .batch-progress__title {
    font-size: var(--font-size-large);
    font-weight: var(--font-weight-bold);
    margin: 0;
  }

  .batch-progress__bar {
    height: 8px;
    background: #E0E0E0;
    border-radius: 4px;
    overflow: hidden;
  }

  .batch-progress__fill {
    height: 100%;
    background: var(--color-primary); /* Blue */
    transition: width 0.3s ease;
  }
  ```
- [ ] Implement `updateBatchProgress(completed, total)` function in `web/js/upload.js`:
  ```javascript
  updateBatchProgress(completed, total) {
    const batchProgress = document.getElementById('batch-progress');
    const title = document.getElementById('batch-progress-title');
    const fill = document.getElementById('batch-progress-fill');

    if (total === 0) {
      batchProgress.setAttribute('hidden', '');
      return;
    }

    batchProgress.removeAttribute('hidden');
    title.textContent = `Converting ${completed} of ${total}...`;
    const percentage = (completed / total) * 100;
    fill.style.width = `${percentage}%`;
  }
  ```

### Task 2: Implement Per-File Status Indicators (AC-2)
- [ ] Update `createFileCard()` in `upload.js` to include status section:
  ```javascript
  createFileCard(file) {
    // ... existing code ...
    card.innerHTML = `
      <div class="file-card__header">
        <span class="file-card__filename" title="${file.name}">
          ${this.truncateFilename(file.name, 30)}
        </span>
        <span class="badge badge--${format.toLowerCase()}">${format}</span>
      </div>
      <div class="file-card__body">
        <div class="file-card__size">${fileSize}</div>
        <div class="file-card__status" data-status="queued">
          <span class="status-icon">⏱️</span>
          <span class="status-text">Queued</span>
        </div>
      </div>
      <div class="file-card__footer">
        <button class="file-card__download" data-file-id="${fileId}" hidden>
          Download ${format}
        </button>
        <button class="file-card__remove" data-file-id="${fileId}">
          <span class="remove-icon">✕</span> Remove
        </button>
      </div>
    `;
    return { card, fileId };
  }
  ```
- [ ] Add status state styles to `web/css/components.css`:
  ```css
  /* Queued state */
  .file-card[data-status="queued"] {
    background: white;
  }
  .file-card[data-status="queued"] .status-icon {
    color: #999;
  }

  /* Processing state */
  .file-card[data-status="processing"] {
    background: #F0F7FF; /* Light blue */
  }
  .file-card[data-status="processing"] .status-icon {
    color: var(--color-primary); /* Blue */
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }

  /* Complete state */
  .file-card[data-status="complete"] {
    background: #E8F5E9; /* Light green */
  }
  .file-card[data-status="complete"] .status-icon {
    color: #4CAF50; /* Green */
  }

  /* Error state */
  .file-card[data-status="error"] {
    background: #FFEBEE; /* Light red */
  }
  .file-card[data-status="error"] .status-icon {
    color: #D32F2F; /* Red */
  }
  .file-card[data-status="error"] .status-text {
    color: #D32F2F;
    font-size: var(--font-size-small);
  }
  ```

### Task 3: Implement updateFileStatus() Function (AC-2, AC-3)
- [ ] Add `updateFileStatus(fileId, status, errorMessage)` to `upload.js`:
  ```javascript
  updateFileStatus(fileId, status, errorMessage = null) {
    const fileData = this.files.get(fileId);
    if (!fileData) return;

    fileData.status = status;
    const card = fileData.element;
    const statusContainer = card.querySelector('.file-card__status');
    const icon = statusContainer.querySelector('.status-icon');
    const text = statusContainer.querySelector('.status-text');
    const downloadBtn = card.querySelector('.file-card__download');

    // Update card data-status attribute (triggers CSS transition)
    card.setAttribute('data-status', status);

    // Update icon and text based on status
    switch (status) {
      case 'queued':
        icon.textContent = '⏱️';
        text.textContent = 'Queued';
        downloadBtn.setAttribute('hidden', '');
        break;
      case 'processing':
        icon.textContent = '⏳';
        text.textContent = 'Converting...';
        downloadBtn.setAttribute('hidden', '');
        break;
      case 'complete':
        icon.textContent = '✓';
        text.textContent = 'Complete';
        downloadBtn.removeAttribute('hidden');
        break;
      case 'error':
        icon.textContent = '✕';
        text.textContent = `Error: ${errorMessage || 'Conversion failed'}`;
        downloadBtn.setAttribute('hidden', '');
        break;
    }
  }
  ```
- [ ] Add CSS transition for smooth state changes:
  ```css
  .file-card {
    transition: background-color 0.3s ease;
  }

  .file-card__status {
    transition: color 0.3s ease;
  }
  ```

### Task 4: Implement Batch Conversion with Progress Updates (AC-1, AC-4)
- [ ] Add `startBatchConversion()` function to `web/js/app.js`:
  ```javascript
  async function startBatchConversion(uploadManager, targetFormat) {
    const files = Array.from(uploadManager.files.values());
    const total = files.length;
    let completed = 0;

    uploadManager.updateBatchProgress(completed, total);

    for (const fileData of files) {
      // Update to processing state
      uploadManager.updateFileStatus(fileData.fileId, 'processing');

      try {
        // Perform conversion (WASM call)
        const result = await convertFile(fileData.file, targetFormat);

        // Update to complete state
        uploadManager.updateFileStatus(fileData.fileId, 'complete');
        fileData.outputData = result.data;
        fileData.outputFormat = targetFormat;

        completed++;
      } catch (error) {
        // Update to error state
        uploadManager.updateFileStatus(fileData.fileId, 'error', error.message);
        completed++;
      }

      uploadManager.updateBatchProgress(completed, total);
    }

    // All files processed
    const successCount = files.filter(f => f.status === 'complete').length;
    const errorCount = files.filter(f => f.status === 'error').length;

    if (errorCount === 0) {
      uploadManager.showBatchComplete(total);
    } else {
      uploadManager.showBatchComplete(successCount, errorCount);
    }
  }
  ```
- [ ] Add `showBatchComplete(successCount, errorCount)` to `upload.js`:
  ```javascript
  showBatchComplete(successCount, errorCount = 0) {
    const title = document.getElementById('batch-progress-title');
    const cancelBtn = document.getElementById('cancel-batch');

    if (errorCount === 0) {
      title.textContent = `All ${successCount} files converted successfully`;
    } else {
      title.textContent = `${successCount} of ${successCount + errorCount} converted, ${errorCount} failed`;
    }

    cancelBtn.setAttribute('hidden', ''); // Hide cancel button
  }
  ```

### Task 5: Implement Download Button Functionality (AC-5)
- [ ] Add download button click handler to `upload.js`:
  ```javascript
  addFileCard(file) {
    const { card, fileId } = this.createFileCard(file);
    this.fileGrid.appendChild(card);
    this.files.set(fileId, { file, element: card, status: 'queued', fileId });

    // Remove button listener
    card.querySelector('.file-card__remove').addEventListener('click', () => {
      this.removeFile(fileId);
    });

    // Download button listener
    card.querySelector('.file-card__download').addEventListener('click', () => {
      this.downloadFile(fileId);
    });
  }

  downloadFile(fileId) {
    const fileData = this.files.get(fileId);
    if (!fileData || !fileData.outputData) return;

    const originalName = fileData.file.name.split('.').slice(0, -1).join('.');
    const extension = this.getFormatExtension(fileData.outputFormat);
    const filename = `${originalName}_converted.${extension}`;

    const blob = new Blob([fileData.outputData], { type: 'application/octet-stream' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    a.click();
    URL.revokeObjectURL(url);
  }

  getFormatExtension(format) {
    const extensions = {
      'NP3': 'np3',
      'XMP': 'xmp',
      'lrtemplate': 'lrtemplate',
      'Capture One': 'costyle',
      'DCP': 'dcp'
    };
    return extensions[format] || 'bin';
  }
  ```

### Task 6: Implement Error Handling with Specific Messages (AC-6)
- [ ] Add error message mapping in `app.js`:
  ```javascript
  async function convertFile(file, targetFormat) {
    try {
      // Read file as ArrayBuffer
      const buffer = await file.arrayBuffer();
      const data = new Uint8Array(buffer);

      // Detect source format (call WASM detectFormat)
      const sourceFormat = detectFormat(data, file.name);
      if (!sourceFormat) {
        throw new Error('Invalid file format');
      }

      // Parse source file (call WASM parser)
      const recipe = parseFormat(data, sourceFormat);
      if (!recipe) {
        throw new Error(`Parsing failed: Unable to parse ${sourceFormat} file`);
      }

      // Generate target file (call WASM generator)
      const output = generateFormat(recipe, targetFormat);
      if (!output) {
        throw new Error(`Conversion failed: Unable to generate ${targetFormat} file`);
      }

      return { data: output, format: targetFormat };
    } catch (error) {
      // Map generic errors to specific messages
      if (error.message.includes('invalid magic')) {
        throw new Error('Invalid file format');
      } else if (error.message.includes('parse error')) {
        throw new Error(`Parsing failed: ${error.message}`);
      } else {
        throw new Error(error.message);
      }
    }
  }
  ```
- [ ] Test error handling:
  - Upload corrupted .np3 file → "Invalid file format"
  - Upload .txt file renamed to .xmp → "Parsing failed: Invalid XML"
  - WASM conversion error → "Conversion failed: <reason>"

### Task 7: Implement Cancel Batch Conversion (AC-7)
- [ ] Add `cancelBatchConversion()` function to `app.js`:
  ```javascript
  let conversionAborted = false;

  async function startBatchConversion(uploadManager, targetFormat) {
    const files = Array.from(uploadManager.files.values());
    const total = files.length;
    let completed = 0;
    conversionAborted = false;

    uploadManager.updateBatchProgress(completed, total);

    for (const fileData of files) {
      if (conversionAborted) {
        // Stop processing remaining files
        break;
      }

      // ... existing conversion logic ...
    }

    if (conversionAborted) {
      uploadManager.showBatchCancelled(completed, total);
    } else {
      uploadManager.showBatchComplete(successCount, errorCount);
    }
  }

  function cancelBatchConversion() {
    const confirmed = confirm('Cancel all in-progress conversions?');
    if (confirmed) {
      conversionAborted = true;
    }
  }
  ```
- [ ] Add cancel button listener in `upload.js`:
  ```javascript
  constructor() {
    // ... existing code ...

    const cancelBtn = document.getElementById('cancel-batch');
    cancelBtn.addEventListener('click', () => {
      cancelBatchConversion();
    });
  }

  showBatchCancelled(completed, total) {
    const title = document.getElementById('batch-progress-title');
    const cancelBtn = document.getElementById('cancel-batch');

    title.textContent = `Conversion cancelled (${completed} of ${total} complete)`;
    cancelBtn.setAttribute('hidden', '');
  }
  ```

### Task 8: Manual Testing
- [ ] Test batch progress indicator:
  - Upload 10 files → Start conversion → Verify progress updates (0/10, 1/10, ..., 10/10)
  - Verify progress bar fills from 0% to 100%
  - Verify "All 10 files converted successfully" message
- [ ] Test per-file status transitions:
  - Queued → Processing → Complete (verify icon, text, background color change)
  - Queued → Processing → Error (verify error icon, message)
  - Verify smooth CSS transitions (no flicker or jank)
- [ ] Test spinner animation:
  - Verify spinner rotates during "processing" state
  - Verify animation stops when conversion completes
- [ ] Test download functionality:
  - Complete conversion → Click "Download XMP" → Verify file downloads
  - Verify filename: `preset_converted.xmp`
  - Verify downloaded file is valid (open in text editor, check content)
- [ ] Test error messages:
  - Upload corrupted file → Verify "Invalid file format"
  - Upload .txt renamed to .xmp → Verify "Parsing failed: ..."
  - Verify error message truncated if >50 chars (full message on hover)
- [ ] Test cancel functionality:
  - Start batch conversion (10 files) → Click "Cancel All" after 3 files
  - Verify confirmation dialog appears
  - Confirm → Verify remaining files stop processing
  - Verify message: "Conversion cancelled (3 of 10 complete)"
- [ ] Test edge cases:
  - Upload 1 file → Verify singular "Converting 1 of 1..."
  - All files fail → Verify "0 of 10 converted, 10 failed"
  - Mixed success/failure → Verify "7 of 10 converted, 3 failed"

## Dev Notes

### Learnings from Previous Story

**From Story 10-2-batch-file-upload (Status: drafted)**

Previous story not yet implemented. Story 10.3 builds on the file upload foundation by adding conversion progress tracking.

**Reuse from Story 10-2:**
- `web/js/upload.js` - UploadManager class (file card creation, Map-based state management)
- `web/css/components.css` - File card styles (.file-card, .file-card__header, etc.)
- `this.files` Map structure: `Map<fileId, FileData>` where FileData includes `{ file, element, status }`
- File card HTML structure (header, body, footer sections)
- Responsive grid layout (3-column → 2-column → 1-column)

**New Status States Added (Story 10.3):**
- Queued → Processing → Complete → Download
- Queued → Processing → Error → Retry (optional)
- Status state stored in `fileData.status` field

**Integration:**
- Add `data-status` attribute to file cards for CSS state-based styling
- Extend `FileData` interface with `outputData`, `outputFormat`, `error` fields
- Batch progress component placed above file grid (Story 10-2 created grid structure)

[Source: docs/stories/10-2-batch-file-upload.md]

### Architecture Alignment

**Tech Spec Epic 10 Alignment:**

Story 10.3 implements **AC-4 (Progress Indicators)** from tech-spec-epic-10.md.

**Conversion Flow with Progress:**
```
Upload files (Story 10-2) → Start batch conversion (Story 10-3)
                                        ↓
          For each file: Queued → Processing → Complete/Error
                                        ↓
                         Update batch progress (3/10, 4/10, ...)
                                        ↓
                     Show download buttons (Complete) or error (Error)
```

**State Management Extension:**
```javascript
interface FileData {
  file: File;              // Original File object (from Story 10-2)
  element: HTMLElement;    // File card DOM element (from Story 10-2)
  fileId: string;          // Unique file ID (from Story 10-2)
  status: 'queued' | 'processing' | 'complete' | 'error';  // NEW
  outputData?: Uint8Array; // Converted file data (NEW)
  outputFormat?: string;   // Target format (e.g., 'XMP') (NEW)
  error?: string;          // Error message if status='error' (NEW)
}
```

**New Functions Added:**
```javascript
// web/js/upload.js
updateFileStatus(fileId, status, errorMessage)   // Update file card UI
updateBatchProgress(completed, total)            // Update progress bar
showBatchComplete(successCount, errorCount)      // Show completion message
showBatchCancelled(completed, total)             // Show cancellation message
downloadFile(fileId)                             // Trigger browser download
getFormatExtension(format)                       // Map format to extension

// web/js/app.js (NEW)
startBatchConversion(uploadManager, targetFormat) // Orchestrate batch conversion
convertFile(file, targetFormat)                  // Convert single file (WASM)
cancelBatchConversion()                          // Abort in-progress batch
```

[Source: docs/tech-spec-epic-10.md#Detailed-Design]

### Batch Conversion Architecture

**Sequential Processing (Not Parallel):**

Recipe uses sequential file processing to maintain simplicity:
```javascript
for (const fileData of files) {
  updateFileStatus(fileId, 'processing');
  const result = await convertFile(file, targetFormat);
  updateFileStatus(fileId, 'complete');
}
```

**Why Sequential?**
- WASM conversions are extremely fast (<100ms per file)
- Sequential processing prevents race conditions
- Simpler state management (one file processing at a time)
- Progress updates are deterministic (no out-of-order completions)

**Performance:**
- 10 files × 100ms = 1 second total (acceptable UX)
- Batch cancellation is immediate (break loop)
- UI updates are synchronous (no Promise.all complexity)

**Future Optimization (if needed):**
- Parallel processing with Web Workers (defer to post-MVP)
- Worker pool (4 workers, process 4 files simultaneously)
- Only if user feedback indicates slow batch processing

[Source: docs/tech-spec-epic-10.md#Performance-Considerations]

### WASM Integration

**Conversion Function Signature:**

```javascript
// web/js/wasm-bridge.js (from Story 2-6)
export async function convertFile(fileData, targetFormat) {
  // 1. Detect source format
  const sourceFormat = await wasmDetectFormat(fileData, filename);

  // 2. Parse source file
  const recipe = await wasmParse(fileData, sourceFormat);

  // 3. Generate target file
  const output = await wasmGenerate(recipe, targetFormat);

  return output; // Uint8Array
}
```

**Error Handling:**

WASM functions return errors via exception messages:
- `"NP3: invalid magic bytes"` → "Invalid file format"
- `"XMP: parse error at line 42"` → "Parsing failed: Invalid XML"
- `"lrtemplate: unsupported version 12"` → "Parsing failed: Unsupported version"

Story 10.3 wraps WASM errors into user-friendly messages.

[Source: docs/stories/2-6-wasm-conversion-execution.md#WASM-Error-Handling]

### CSS State-Based Styling

**Data Attribute Pattern:**

Recipe uses `data-status` attribute for state-based styling:
```html
<div class="file-card" data-status="processing">
  <!-- Card content -->
</div>
```

**CSS Selectors:**
```css
.file-card[data-status="queued"] { background: white; }
.file-card[data-status="processing"] { background: #F0F7FF; }
.file-card[data-status="complete"] { background: #E8F5E9; }
.file-card[data-status="error"] { background: #FFEBEE; }
```

**Why Data Attributes?**
- CSS transitions apply automatically when attribute changes
- JavaScript only updates one attribute: `card.setAttribute('data-status', 'complete')`
- No manual class manipulation (`add`/`remove` multiple classes)
- State is self-documenting in HTML (inspect element shows current state)

[Source: CSS Best Practices - Data Attributes for State Management]

### Download Functionality

**Blob API for File Downloads:**

```javascript
downloadFile(fileId) {
  const fileData = this.files.get(fileId);

  // Create Blob from Uint8Array
  const blob = new Blob([fileData.outputData], { type: 'application/octet-stream' });

  // Create temporary URL
  const url = URL.createObjectURL(blob);

  // Create <a> element and trigger download
  const a = document.createElement('a');
  a.href = url;
  a.download = 'preset_converted.xmp';
  a.click();

  // Clean up temporary URL
  URL.revokeObjectURL(url);
}
```

**Filename Convention:**
- Original: `my-preset.np3`
- Downloaded: `my-preset_converted.xmp` (original name + `_converted` + new extension)

**Browser Compatibility:**
- Chrome, Firefox, Edge: Full support
- Safari: Full support (iOS 13+)
- Blob API universally supported (no polyfill needed)

[Source: MDN Web Docs - Blob API]

### Progress Bar Animation

**CSS Transition for Smooth Progress:**

```css
.batch-progress__fill {
  width: 0%; /* Initial state */
  transition: width 0.3s ease; /* Smooth animation */
}
```

**JavaScript Width Updates:**
```javascript
const percentage = (completed / total) * 100;
fill.style.width = `${percentage}%`; // Triggers CSS transition
```

**Visual Effect:**
- Progress bar smoothly expands from 0% → 10% → 20% → ... → 100%
- Each update takes 0.3 seconds (ease timing function)
- No janky jumps (CSS handles interpolation)

**Performance:**
- CSS transitions are GPU-accelerated (smooth on low-end devices)
- Width updates are batched by browser (no reflow storm)

[Source: CSS Transitions Best Practices]

### Project Structure Notes

**Modified Files (Story 10.3):**
- `web/index.html` - Add batch progress component (above file grid)
- `web/css/components.css` - Add batch progress styles, status state styles
- `web/js/upload.js` - Extend UploadManager with status updates, download functionality
- `web/js/app.js` - Add batch conversion orchestration (startBatchConversion, convertFile)

**No New Files Created:**
- Story 10.3 extends existing modules from Story 10-2 and Story 2-6 (WASM bridge)

**Files from Previous Stories (Reused):**
- `web/js/upload.js` - UploadManager class (Story 10-2)
- `web/js/wasm-bridge.js` - WASM conversion functions (Story 2-6)
- `web/css/components.css` - File card styles (Story 10-2)

[Source: docs/tech-spec-epic-10.md#Services-and-Modules]

### Testing Strategy

**Manual Testing (Required):**

1. **Batch Progress Testing:**
   - Upload 10 files → Start conversion → Verify progress: 0/10, 1/10, ..., 10/10
   - Verify progress bar fills from 0% → 100%
   - Verify completion message: "All 10 files converted successfully"

2. **Status Transition Testing:**
   - Verify Queued → Processing → Complete (icon, text, background color)
   - Verify Queued → Processing → Error (error icon, message)
   - Verify smooth CSS transitions (no flicker)

3. **Download Testing:**
   - Complete conversion → Click "Download XMP" → Verify file downloads
   - Verify filename format: `<name>_converted.xmp`
   - Open downloaded file → Verify valid XML/binary content

4. **Error Handling Testing:**
   - Upload corrupted .np3 → Verify "Invalid file format"
   - Upload .txt renamed to .xmp → Verify "Parsing failed: ..."
   - Verify error message truncation (>50 chars)

5. **Cancel Testing:**
   - Start batch (10 files) → Cancel after 3 files
   - Verify confirmation dialog
   - Verify remaining files stop processing
   - Verify message: "Conversion cancelled (3 of 10 complete)"

**Automated Testing (Optional):**

Unit tests for status update logic:
```javascript
// test/upload.test.js
test('updateFileStatus updates card UI correctly', () => {
  const upload = new UploadManager();
  const fileId = 'test-file-123';

  upload.updateFileStatus(fileId, 'processing');
  expect(card.getAttribute('data-status')).toBe('processing');
  expect(icon.textContent).toBe('⏳');

  upload.updateFileStatus(fileId, 'complete');
  expect(card.getAttribute('data-status')).toBe('complete');
  expect(icon.textContent).toBe('✓');
});
```

[Source: docs/tech-spec-epic-10.md#Test-Strategy-Summary]

### Known Risks

**RISK-32: Conversion errors not handled gracefully**
- **Impact**: User sees generic "Error" without actionable information
- **Mitigation**: Map WASM errors to specific user-friendly messages
- **Test**: Upload corrupted files, verify error messages are clear

**RISK-33: Batch progress UI blocks user interaction**
- **Impact**: User can't remove files or cancel during batch processing
- **Mitigation**: Ensure remove buttons work during processing, add "Cancel All" button
- **Acceptable**: User can cancel anytime, remove completed files

**RISK-34: Large batch (100+ files) freezes UI**
- **Impact**: Browser becomes unresponsive during long batch conversion
- **Mitigation**: Sequential processing with `await` prevents UI blocking (WASM runs off main thread)
- **Performance**: 100 files × 100ms = 10 seconds (acceptable for web UI)

[Source: docs/tech-spec-epic-10.md#Risks-Assumptions-Open-Questions]

### References

- [Source: docs/tech-spec-epic-10.md#Acceptance-Criteria] - AC-4: Progress Indicators
- [Source: docs/tech-spec-epic-10.md#Detailed-Design] - Batch conversion flow, state management
- [Source: docs/stories/10-2-batch-file-upload.md] - UploadManager class, file card structure
- [Source: docs/stories/2-6-wasm-conversion-execution.md] - WASM conversion API, error handling
- [MDN: Blob API](https://developer.mozilla.org/en-US/docs/Web/API/Blob)
- [MDN: URL.createObjectURL](https://developer.mozilla.org/en-US/docs/Web/API/URL/createObjectURL)
- [CSS Transitions](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Transitions)

## Dev Agent Record

### Context Reference

- docs/stories/10-3-progress-indicators.context.xml

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

### Completion Notes List

### File List
