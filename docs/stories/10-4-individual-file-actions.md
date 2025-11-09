# Story 10.4: Individual File Conversion Controls

Status: ready-for-dev

## Story

As a **photographer with files in different formats**,
I want **per-file conversion controls to select different target formats for each file**,
so that **I can convert files to different formats in a single batch without running separate conversions**.

## Acceptance Criteria

**AC-1: Per-File Format Selection Dropdown**
- ✅ Each file card displays "Convert to..." dropdown
- ✅ Dropdown shows only valid target formats for the source format:
  - NP3 → XMP, lrtemplate, .costyle, DCP (4 options)
  - XMP → NP3, lrtemplate, .costyle, DCP (4 options)
  - lrtemplate → NP3, XMP, .costyle, DCP (4 options)
  - .costyle → NP3, XMP, lrtemplate, DCP (4 options)
  - DCP → NP3, XMP, lrtemplate, .costyle (4 options)
- ✅ Dropdown disabled during conversion ("processing" state)
- ✅ Dropdown hidden when conversion complete ("complete" state)
- ✅ Default selection: First valid target format (alphabetically)

**AC-2: Individual "Convert" Button**
- ✅ Each file card has "Convert" button
- ✅ Button enabled when file in "queued" state
- ✅ Button disabled during conversion ("processing" state)
- ✅ Button hidden when conversion complete ("complete" state, replaced by download)
- ✅ Button text: "Convert" (no icon needed)
- ✅ Clicking button starts conversion for that file only (not batch)

**AC-3: Individual File Conversion Flow**
- ✅ User selects target format from dropdown
- ✅ User clicks "Convert" button
- ✅ File status changes: queued → processing → complete/error
- ✅ Only selected file converts (other files remain in queued state)
- ✅ User can convert files in any order (file 3 before file 1, etc.)
- ✅ Batch progress indicator shows per-file progress (not required to convert all files)

**AC-4: Mixed Format Conversion Support**
- ✅ User can select different target formats for different files:
  - File 1 (NP3) → Convert to XMP
  - File 2 (XMP) → Convert to lrtemplate
  - File 3 (lrtemplate) → Convert to DCP
- ✅ Each file converts to its selected target format
- ✅ Download buttons show correct format: "Download XMP", "Download lrtemplate", etc.
- ✅ Downloaded filenames reflect target format: `preset_converted.xmp`, `style_converted.lrtemplate`

**AC-5: Retry Failed Conversions**
- ✅ Files in "error" state show "Retry" button (replaces "Convert" button)
- ✅ Clicking "Retry" re-attempts conversion with same target format
- ✅ Error message clears when retry starts (status → processing)
- ✅ If retry succeeds → status changes to complete
- ✅ If retry fails → error message updates with new error (may differ from first attempt)
- ✅ Retry attempts not limited (user can retry indefinitely)

**AC-6: Batch vs Individual Conversion Coexistence**
- ✅ Batch "Convert All" button (Story 10-5) converts all queued files to same format
- ✅ Individual "Convert" buttons convert single files to custom formats
- ✅ User can mix: batch convert 5 files, then individually convert remaining 2
- ✅ Batch progress tracks both batch-converted and individually-converted files
- ✅ No conflict: user can't trigger batch and individual conversions simultaneously

**AC-7: Conversion Cancel for Individual Files**
- ✅ Individual file conversion can be cancelled:
  - During "processing" state, file card shows "Cancel" button (replaces "Convert")
  - Clicking "Cancel" aborts conversion for that file
  - File status reverts to "queued" (user can convert again)
- ✅ Cancelling individual file doesn't affect other files
- ✅ No confirmation dialog for individual cancel (only batch cancel requires confirmation)

## Tasks / Subtasks

### Task 1: Add Format Selection Dropdown to File Cards (AC-1)
- [ ] Update `createFileCard()` in `upload.js` to include format dropdown:
  ```javascript
  createFileCard(file) {
    const fileId = `file-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
    const sourceFormat = this.detectFormat(file.name);
    const targetFormats = this.getValidTargetFormats(sourceFormat);
    const fileSize = this.formatFileSize(file.size);

    const card = document.createElement('div');
    card.className = 'file-card';
    card.id = fileId;
    card.setAttribute('data-status', 'queued');
    card.innerHTML = `
      <div class="file-card__header">
        <span class="file-card__filename" title="${file.name}">
          ${this.truncateFilename(file.name, 30)}
        </span>
        <span class="badge badge--${sourceFormat.toLowerCase()}">${sourceFormat}</span>
      </div>
      <div class="file-card__body">
        <div class="file-card__size">${fileSize}</div>
        <div class="file-card__status">
          <span class="status-icon">⏱️</span>
          <span class="status-text">Queued</span>
        </div>
      </div>
      <div class="file-card__conversion">
        <label class="file-card__label">Convert to:</label>
        <select class="file-card__format-select" data-file-id="${fileId}">
          ${targetFormats.map(fmt => `<option value="${fmt}">${fmt}</option>`).join('')}
        </select>
        <button class="file-card__convert button button--primary" data-file-id="${fileId}">
          Convert
        </button>
      </div>
      <div class="file-card__footer">
        <button class="file-card__download button button--primary" data-file-id="${fileId}" hidden>
          Download ${targetFormats[0]}
        </button>
        <button class="file-card__remove button button--secondary" data-file-id="${fileId}">
          Remove
        </button>
      </div>
    `;

    return { card, fileId, sourceFormat };
  }

  getValidTargetFormats(sourceFormat) {
    const allFormats = ['NP3', 'XMP', 'lrtemplate', 'Capture One', 'DCP'];
    // Return all formats except source format
    return allFormats.filter(fmt => fmt !== sourceFormat).sort();
  }
  ```

### Task 2: Style Format Dropdown and Convert Button (AC-1, AC-2)
- [ ] Add file card conversion section styles to `web/css/components.css`:
  ```css
  .file-card__conversion {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin: 1rem 0;
    padding-top: 1rem;
    border-top: 1px solid var(--color-border);
  }

  .file-card__label {
    font-size: var(--font-size-small);
    font-weight: var(--font-weight-bold);
    color: var(--color-text-secondary);
  }

  .file-card__format-select {
    flex: 1;
    padding: 0.5rem;
    border: 1px solid var(--color-border);
    border-radius: 4px;
    background: white;
    font-size: var(--font-size-base);
    cursor: pointer;
  }

  .file-card__format-select:disabled {
    background: #f5f5f5;
    cursor: not-allowed;
    opacity: 0.6;
  }

  .file-card__convert {
    padding: 0.5rem 1rem;
    white-space: nowrap;
  }

  /* Hide conversion section when complete */
  .file-card[data-status="complete"] .file-card__conversion,
  .file-card[data-status="error"] .file-card__conversion {
    display: none;
  }

  /* Show retry button for error state */
  .file-card[data-status="error"] .file-card__conversion {
    display: flex; /* Re-enable for retry */
  }

  .file-card[data-status="error"] .file-card__convert::before {
    content: "Retry ";
  }
  ```

### Task 3: Implement Individual File Conversion (AC-2, AC-3)
- [ ] Add individual conversion handler to `upload.js`:
  ```javascript
  addFileCard(file) {
    const { card, fileId, sourceFormat } = this.createFileCard(file);
    this.fileGrid.appendChild(card);
    this.files.set(fileId, {
      file,
      element: card,
      status: 'queued',
      fileId,
      sourceFormat,
      targetFormat: null,
      outputData: null
    });

    // Convert button listener
    card.querySelector('.file-card__convert').addEventListener('click', () => {
      this.convertIndividualFile(fileId);
    });

    // Format select listener (update target format)
    card.querySelector('.file-card__format-select').addEventListener('change', (e) => {
      const fileData = this.files.get(fileId);
      fileData.targetFormat = e.target.value;
    });

    // Download button listener
    card.querySelector('.file-card__download').addEventListener('click', () => {
      this.downloadFile(fileId);
    });

    // Remove button listener
    card.querySelector('.file-card__remove').addEventListener('click', () => {
      this.removeFile(fileId);
    });

    // Set initial target format (first option in dropdown)
    const select = card.querySelector('.file-card__format-select');
    const fileData = this.files.get(fileId);
    fileData.targetFormat = select.value;
  }

  async convertIndividualFile(fileId) {
    const fileData = this.files.get(fileId);
    if (!fileData || fileData.status !== 'queued') return;

    const targetFormat = fileData.targetFormat;
    if (!targetFormat) {
      console.error('No target format selected');
      return;
    }

    // Update to processing state
    this.updateFileStatus(fileId, 'processing');

    // Disable dropdown and convert button during processing
    const card = fileData.element;
    const select = card.querySelector('.file-card__format-select');
    const convertBtn = card.querySelector('.file-card__convert');
    select.disabled = true;
    convertBtn.disabled = true;

    try {
      // Call conversion function (WASM bridge)
      const result = await convertFile(fileData.file, targetFormat);

      // Update to complete state
      this.updateFileStatus(fileId, 'complete');
      fileData.outputData = result.data;
      fileData.outputFormat = targetFormat;

      // Update download button text
      const downloadBtn = card.querySelector('.file-card__download');
      downloadBtn.textContent = `Download ${targetFormat}`;
      downloadBtn.removeAttribute('hidden');

    } catch (error) {
      // Update to error state
      this.updateFileStatus(fileId, 'error', error.message);

      // Re-enable dropdown and convert button (for retry)
      select.disabled = false;
      convertBtn.disabled = false;
    }
  }
  ```

### Task 4: Implement Mixed Format Conversion Support (AC-4)
- [ ] Update `downloadFile()` to use per-file target format:
  ```javascript
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
- [ ] Test mixed format conversion:
  - Upload 3 files (NP3, XMP, lrtemplate)
  - Convert NP3 → XMP
  - Convert XMP → DCP
  - Convert lrtemplate → NP3
  - Verify each downloads with correct format

### Task 5: Implement Retry Failed Conversions (AC-5)
- [ ] Update `updateFileStatus()` to handle retry state:
  ```javascript
  updateFileStatus(fileId, status, errorMessage = null) {
    const fileData = this.files.get(fileId);
    if (!fileData) return;

    fileData.status = status;
    const card = fileData.element;
    const statusContainer = card.querySelector('.file-card__status');
    const icon = statusContainer.querySelector('.status-icon');
    const text = statusContainer.querySelector('.status-text');
    const convertBtn = card.querySelector('.file-card__convert');
    const select = card.querySelector('.file-card__format-select');

    // Update card data-status attribute
    card.setAttribute('data-status', status);

    switch (status) {
      case 'queued':
        icon.textContent = '⏱️';
        text.textContent = 'Queued';
        select.disabled = false;
        convertBtn.disabled = false;
        convertBtn.textContent = 'Convert';
        break;
      case 'processing':
        icon.textContent = '⏳';
        text.textContent = 'Converting...';
        select.disabled = true;
        convertBtn.disabled = true;
        break;
      case 'complete':
        icon.textContent = '✓';
        text.textContent = 'Complete';
        break;
      case 'error':
        icon.textContent = '✕';
        text.textContent = `Error: ${errorMessage || 'Conversion failed'}`;
        select.disabled = false;
        convertBtn.disabled = false;
        convertBtn.textContent = 'Retry'; // Change button text to Retry
        fileData.error = errorMessage;
        break;
    }
  }
  ```
- [ ] Test retry functionality:
  - Upload corrupted file → Conversion fails → "Retry" button appears
  - Click "Retry" → Conversion re-attempts
  - Verify error clears when retry starts

### Task 6: Implement Batch and Individual Conversion Coexistence (AC-6)
- [ ] Add batch "Convert All" functionality in `web/js/app.js`:
  ```javascript
  // Global batch conversion
  async function startBatchConversion(uploadManager) {
    const batchFormatSelect = document.getElementById('batch-format-select');
    const batchTargetFormat = batchFormatSelect.value;

    const files = Array.from(uploadManager.files.values());
    const queuedFiles = files.filter(f => f.status === 'queued');

    if (queuedFiles.length === 0) {
      alert('No files to convert');
      return;
    }

    for (const fileData of queuedFiles) {
      // Set target format to batch format
      fileData.targetFormat = batchTargetFormat;

      // Convert individual file
      await uploadManager.convertIndividualFile(fileData.fileId);
    }
  }
  ```
- [ ] Add batch format selector to `web/index.html`:
  ```html
  <div class="batch-controls">
    <label for="batch-format-select">Convert all to:</label>
    <select id="batch-format-select" class="batch-format-select">
      <option value="NP3">NP3</option>
      <option value="XMP">XMP</option>
      <option value="lrtemplate">lrtemplate</option>
      <option value="Capture One">Capture One</option>
      <option value="DCP">DCP</option>
    </select>
    <button id="convert-all-btn" class="button button--primary">Convert All</button>
  </div>
  ```
- [ ] Test batch/individual coexistence:
  - Upload 10 files
  - Batch convert 5 files to XMP
  - Individually convert 2 files to DCP
  - Individually convert 1 file to lrtemplate
  - Verify batch progress shows "8 of 10 converted"

### Task 7: Implement Individual File Cancellation (AC-7)
- [ ] Add cancellation support to `convertIndividualFile()`:
  ```javascript
  async convertIndividualFile(fileId) {
    const fileData = this.files.get(fileId);
    if (!fileData || fileData.status !== 'queued') return;

    const targetFormat = fileData.targetFormat;
    const card = fileData.element;
    const select = card.querySelector('.file-card__format-select');
    const convertBtn = card.querySelector('.file-card__convert');

    // Update to processing state
    this.updateFileStatus(fileId, 'processing');
    select.disabled = true;

    // Change button to "Cancel"
    convertBtn.textContent = 'Cancel';
    convertBtn.disabled = false;

    // Create AbortController for cancellation
    const abortController = new AbortController();
    fileData.abortController = abortController;

    // Cancel button listener
    const cancelHandler = () => {
      abortController.abort();
      this.updateFileStatus(fileId, 'queued');
      convertBtn.removeEventListener('click', cancelHandler);
    };

    convertBtn.addEventListener('click', cancelHandler);

    try {
      const result = await convertFile(fileData.file, targetFormat, abortController.signal);

      // Check if aborted
      if (abortController.signal.aborted) {
        return; // Conversion cancelled
      }

      // Update to complete state
      this.updateFileStatus(fileId, 'complete');
      fileData.outputData = result.data;
      fileData.outputFormat = targetFormat;

      const downloadBtn = card.querySelector('.file-card__download');
      downloadBtn.textContent = `Download ${targetFormat}`;
      downloadBtn.removeAttribute('hidden');

    } catch (error) {
      if (error.name === 'AbortError') {
        // Conversion was cancelled
        this.updateFileStatus(fileId, 'queued');
      } else {
        // Conversion failed
        this.updateFileStatus(fileId, 'error', error.message);
        select.disabled = false;
        convertBtn.disabled = false;
      }
    } finally {
      convertBtn.removeEventListener('click', cancelHandler);
      delete fileData.abortController;
    }
  }
  ```
- [ ] Test cancellation:
  - Start individual conversion → Click "Cancel" mid-conversion
  - Verify status reverts to "queued"
  - Verify user can convert again after cancelling

### Task 8: Manual Testing
- [ ] Test per-file format selection:
  - Upload NP3 file → Verify dropdown shows XMP, lrtemplate, Capture One, DCP
  - Upload XMP file → Verify dropdown shows NP3, lrtemplate, Capture One, DCP
  - Verify source format excluded from dropdown
- [ ] Test individual conversion:
  - Upload 3 files → Convert file 2 first → Verify only file 2 processes
  - Convert file 1 → Convert file 3 → Verify sequential processing
- [ ] Test mixed format conversion:
  - Upload 5 files → Select different target formats for each
  - Convert all → Verify each downloads with correct format
  - Verify filenames: `preset_converted.xmp`, `style_converted.dcp`
- [ ] Test retry functionality:
  - Upload corrupted file → Verify "Error" state with "Retry" button
  - Click "Retry" → Verify conversion re-attempts
  - Retry multiple times → Verify no limit on retry attempts
- [ ] Test batch/individual coexistence:
  - Upload 10 files → Batch convert 5 to XMP
  - Individually convert 2 to DCP → Verify both work together
  - Verify batch progress: "7 of 10 converted"
- [ ] Test individual cancellation:
  - Start conversion → Click "Cancel" during processing
  - Verify status reverts to "queued"
  - Verify "Convert" button reappears
- [ ] Test edge cases:
  - User changes format dropdown mid-conversion → Verify no effect (dropdown disabled)
  - User removes file during conversion → Verify graceful cleanup
  - User downloads file, then removes it → Verify no errors

## Dev Notes

### Learnings from Previous Story

**From Story 10-3-progress-indicators (Status: drafted)**

Previous story not yet implemented. Story 10.4 builds on the progress tracking foundation by adding per-file conversion controls.

**Reuse from Story 10-3:**
- `updateFileStatus(fileId, status, errorMessage)` - Status update logic
- `FileData` interface with `status`, `outputData`, `outputFormat` fields
- Status states: queued, processing, complete, error
- Download button functionality (`downloadFile()`)
- Error message display in file cards

**New Functionality (Story 10.4):**
- Per-file format selection dropdown (custom target format per file)
- Individual "Convert" button (convert single file, not batch)
- Retry button for error state (re-attempt conversion)
- Individual cancellation (cancel single file, not batch)
- Mixed format conversion (each file converts to different format)

**Integration:**
- Story 10.3 tracks progress for all conversions (batch + individual)
- Story 10.4 adds individual conversion triggers
- Batch progress indicator updates for both batch and individual conversions

[Source: docs/stories/10-3-progress-indicators.md]

### Architecture Alignment

**Tech Spec Epic 10 Alignment:**

Story 10.4 implements **AC-7 (Individual File Conversion)** from tech-spec-epic-10.md.

**Conversion Flow:**
```
User uploads files (Story 10-2)
        ↓
User selects target format per file (Story 10.4)
        ↓
User clicks "Convert" on individual file (Story 10.4)
        ↓
File converts: queued → processing → complete (Story 10.3)
        ↓
User downloads converted file (Story 10.3)
```

**Batch vs Individual Conversion:**
```javascript
// Batch conversion (Story 10-5)
batchConvert(targetFormat) {
  queuedFiles.forEach(file => {
    file.targetFormat = targetFormat; // Same format for all
    convertIndividualFile(file.fileId);
  });
}

// Individual conversion (Story 10.4)
convertIndividualFile(fileId) {
  const file = files.get(fileId);
  const targetFormat = file.targetFormat; // Custom format per file
  await convertFile(file.file, targetFormat);
}
```

**Key Difference:**
- Batch: Set all files to same target format, convert all
- Individual: Each file has custom target format, convert one at a time

[Source: docs/tech-spec-epic-10.md#Detailed-Design]

### Format Validation Logic

**Valid Target Formats Per Source:**

Recipe supports bidirectional conversion between all formats, so each format can convert to all others except itself:

```javascript
getValidTargetFormats(sourceFormat) {
  const allFormats = ['NP3', 'XMP', 'lrtemplate', 'Capture One', 'DCP'];
  return allFormats.filter(fmt => fmt !== sourceFormat).sort();
}
```

**Example:**
- NP3 → [Capture One, DCP, lrtemplate, XMP] (alphabetically sorted)
- XMP → [Capture One, DCP, lrtemplate, NP3]
- DCP → [Capture One, lrtemplate, NP3, XMP]

**No Invalid Conversions:**
All source→target combinations are valid (verified in Epic 1 parameter mapping).

[Source: docs/stories/1-8-parameter-mapping-rules.md#Format-Compatibility-Matrix]

### AbortController for Cancellation

**Modern Cancellation Pattern:**

Recipe uses `AbortController` for async operation cancellation:

```javascript
// Create abort controller
const abortController = new AbortController();

// Pass signal to async function
await convertFile(file, format, abortController.signal);

// Cancel operation
abortController.abort();

// Check if aborted
if (abortController.signal.aborted) {
  return; // Operation cancelled
}
```

**WASM Integration:**

WASM conversion functions check `signal.aborted` between steps:
```javascript
async function convertFile(file, targetFormat, signal) {
  const data = await file.arrayBuffer();

  if (signal?.aborted) throw new DOMException('Aborted', 'AbortError');

  const sourceFormat = detectFormat(data, file.name);

  if (signal?.aborted) throw new DOMException('Aborted', 'AbortError');

  const recipe = parseFormat(data, sourceFormat);

  if (signal?.aborted) throw new DOMException('Aborted', 'AbortError');

  const output = generateFormat(recipe, targetFormat);

  return { data: output, format: targetFormat };
}
```

**Browser Support:**
- Chrome 66+, Firefox 57+, Safari 12.1+, Edge 79+
- No polyfill needed (universally supported)

[Source: MDN Web Docs - AbortController]

### Retry Strategy

**Automatic Retry vs Manual Retry:**

Recipe uses **manual retry** (user clicks "Retry" button):
- No automatic retry (avoids infinite loops on persistent errors)
- User decides when to retry (after fixing file, checking network, etc.)
- No retry limit (user can retry indefinitely if they choose)

**Retry State Machine:**
```
queued → processing → error (conversion fails)
                         ↓
         User clicks "Retry" button
                         ↓
              error → processing (retry attempt)
                         ↓
            complete (success) or error (failure)
```

**Retry Use Cases:**
- Corrupted file → User fixes file → Retry
- Network error during WASM load → Retry
- Temporary browser memory issue → Retry
- User selected wrong target format → Change format → Retry

[Source: UX Best Practices - Error Recovery Patterns]

### Batch and Individual Conversion Coordination

**Shared Conversion Logic:**

Both batch and individual conversions call the same `convertIndividualFile()` function:

```javascript
// Individual conversion (Story 10.4)
convertBtn.addEventListener('click', () => {
  uploadManager.convertIndividualFile(fileId);
});

// Batch conversion (Story 10-5)
async function batchConvert(targetFormat) {
  for (const fileData of queuedFiles) {
    fileData.targetFormat = targetFormat; // Set batch format
    await uploadManager.convertIndividualFile(fileData.fileId); // Reuse same function
  }
}
```

**Benefits:**
- No code duplication (single conversion implementation)
- Consistent error handling (same try/catch logic)
- Unified progress tracking (same `updateFileStatus()` calls)
- Simplified testing (test one function, covers both flows)

**Concurrency Control:**

Recipe prevents simultaneous batch + individual conversions:
```javascript
let conversionInProgress = false;

async function convertIndividualFile(fileId) {
  if (conversionInProgress) {
    alert('Conversion already in progress');
    return;
  }

  conversionInProgress = true;
  try {
    await doConversion();
  } finally {
    conversionInProgress = false;
  }
}
```

[Source: docs/tech-spec-epic-10.md#Concurrency-Model]

### Project Structure Notes

**Modified Files (Story 10.4):**
- `web/index.html` - Add batch format selector, "Convert All" button
- `web/css/components.css` - Add conversion section styles (.file-card__conversion)
- `web/js/upload.js` - Extend UploadManager with `convertIndividualFile()`, retry logic
- `web/js/app.js` - Add batch conversion orchestration (`startBatchConversion()`)

**No New Files Created:**
- Story 10.4 extends existing modules from Story 10-2 and 10-3

**Files from Previous Stories (Reused):**
- `web/js/upload.js` - UploadManager, file cards (Story 10-2)
- `web/js/upload.js` - updateFileStatus, downloadFile (Story 10-3)
- `web/js/wasm-bridge.js` - convertFile function (Story 2-6)

[Source: docs/tech-spec-epic-10.md#Services-and-Modules]

### Testing Strategy

**Manual Testing (Required):**

1. **Format Selection Testing:**
   - Upload each format (NP3, XMP, lrtemplate, .costyle, DCP)
   - Verify dropdown excludes source format
   - Verify dropdown shows 4 target formats (alphabetically)

2. **Individual Conversion Testing:**
   - Upload 5 files → Convert file 3 first (not sequential)
   - Verify only file 3 processes
   - Convert file 1 → file 5 → file 2 (random order)
   - Verify all convert successfully

3. **Mixed Format Conversion Testing:**
   - Upload 3 NP3 files
   - Convert to: XMP, DCP, lrtemplate (different formats)
   - Verify each downloads with correct format
   - Verify filenames: `preset1_converted.xmp`, `preset2_converted.dcp`, etc.

4. **Retry Testing:**
   - Upload corrupted file → Verify "Retry" button
   - Click "Retry" → Verify conversion re-attempts
   - Retry 5 times → Verify no limit

5. **Batch/Individual Coexistence:**
   - Upload 10 files → Batch convert 5 to XMP
   - Individually convert 3 to DCP
   - Verify batch progress: "8 of 10 converted"

6. **Cancellation Testing:**
   - Start conversion → Click "Cancel" mid-conversion
   - Verify status reverts to "queued"
   - Verify "Convert" button reappears

**Automated Testing (Optional):**

Unit tests for format validation:
```javascript
test('getValidTargetFormats excludes source format', () => {
  const upload = new UploadManager();
  const targets = upload.getValidTargetFormats('NP3');
  expect(targets).not.toContain('NP3');
  expect(targets.length).toBe(4); // 5 total formats - 1 source = 4 targets
});

test('getValidTargetFormats returns sorted array', () => {
  const upload = new UploadManager();
  const targets = upload.getValidTargetFormats('XMP');
  expect(targets).toEqual(['Capture One', 'DCP', 'lrtemplate', 'NP3']); // Alphabetical
});
```

[Source: docs/tech-spec-epic-10.md#Test-Strategy-Summary]

### Known Risks

**RISK-35: User confusion with batch vs individual conversion**
- **Impact**: User expects "Convert All" to override individual selections
- **Mitigation**: Clear UI separation (batch controls at top, individual in cards)
- **Documentation**: Add tooltip: "Batch converts all files to same format"

**RISK-36: AbortController not supported in old browsers**
- **Impact**: Cancellation doesn't work in IE11, old Safari
- **Mitigation**: Feature detection, graceful degradation (hide cancel button if unsupported)
- **Acceptable**: Cancel is optional feature (conversions are fast, <100ms)

**RISK-37: Multiple conversions cause memory pressure**
- **Impact**: Browser tab crashes with 100+ files converted simultaneously
- **Mitigation**: Sequential processing (one file at a time) prevents memory spikes
- **Performance**: 100 files × 100ms = 10 seconds (acceptable)

[Source: docs/tech-spec-epic-10.md#Risks-Assumptions-Open-Questions]

### References

- [Source: docs/tech-spec-epic-10.md#Acceptance-Criteria] - AC-7: Individual File Conversion
- [Source: docs/tech-spec-epic-10.md#Detailed-Design] - Conversion flow, batch vs individual
- [Source: docs/stories/10-2-batch-file-upload.md] - UploadManager class, file cards
- [Source: docs/stories/10-3-progress-indicators.md] - Status updates, download functionality
- [Source: docs/stories/1-8-parameter-mapping-rules.md] - Format compatibility matrix
- [Source: docs/stories/2-6-wasm-conversion-execution.md] - WASM conversion API
- [MDN: AbortController](https://developer.mozilla.org/en-US/docs/Web/API/AbortController)
- [MDN: AbortSignal](https://developer.mozilla.org/en-US/docs/Web/API/AbortSignal)

## Dev Agent Record

### Context Reference

- `docs/stories/10-4-individual-file-actions.context.xml` (Generated: 2025-11-09)

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

### Completion Notes List

### File List
