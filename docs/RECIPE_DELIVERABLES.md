Recipe Deliverables (Current)

Profiles and Presets
- Custom DCP (matrix only): output/Nikon_Z_f_Camera_Custom_Recipe.dcp:1
  * Install: docs/USAGE_CUSTOM_DCP.md:1
- HSL preset (hue corrections): output/Recipe_HSL_Adjustments.xmp:1
  * Drop into Lightroom/ACR presets; complements the DCP for closer NX match.

Calibration Artifacts
- Derived matrix: output/nikon_color_matrix.json:1 (camera→XYZ with metrics)
- NX patch targets: output/nx_patch_targets.json:1
- Chroma LUT: output/chroma_lut.json:1
- Hue LUT: output/hue_lut.json:1
- Polar LUT: output/polar_lut.json:1

Evaluation
- Residuals (patch ΔE): output/pipeline_residuals.json:1
  * Baseline mean ≈ 16.46; +chroma ≈ 14.84; +hue ≈ 14.75
  * Polar LUT preserves patches; helps full images visually more than ΔE on 24 patches.

Reverse Engineering Notes
- Bank references: output/polaris_bank_refs.json:1
- Bank region (float view): output/polaris_bank_region.json:1
- Dispatcher call targets: output/dispatch_call_targets.json:1
- Candidate docs: docs/reverse_engineering/polaris_matrix_candidates.md:1, docs/reverse_engineering/polaris_candidate_B_walkthrough.md:1

Usage Flow
1) Install the DCP; optionally also import the XMP HSL preset.
2) For batch output from NEFs with the calibration matrix + LUTs, use:
   `.venv/bin/python scripts/apply_color_matrix.py`
3) To re-derive with fresh NX export, re-run:
   `.venv/bin/python scripts/sample_patch_targets.py output/piccon21_calibration_full.TIF --output output/nx_patch_targets.json`
   `.venv/bin/python scripts/derive_nikon_matrix.py --target-json output/nx_patch_targets.json --lambda 0.2 --gray-weight 8.0`
   Then regenerate previews and residuals.

