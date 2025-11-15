Custom Nikon Z f DCP (Recipe)

Artifacts
- Derived matrix: output/nikon_color_matrix.json:1
- Generated profile: output/Nikon_Z_f_Camera_Custom_Recipe.dcp:1

Install in Lightroom / Camera Raw
1. Preferred (camera‑scoped folder):
   - Windows: `C:\Users\<you>\AppData\Roaming\Adobe\CameraRaw\CameraProfiles\Camera\Nikon Z f\`
   - macOS: `~/Library/Application Support/Adobe/CameraRaw/CameraProfiles/Camera/Nikon Z f/`
   Or run: `python3 scripts/install_dcp_to_user_camera_folder.py output/Nikon_Z_f_Camera_Custom_Recipe.dcp`
2. Restart Lightroom / Camera Raw.
3. In the Profile Browser, ensure “Show Partially Compatible Profiles” is enabled, then select the unique profile name you installed.
4. If a profile “imports” but doesn’t appear, verify it has a unique ProfileName and that it’s installed under the Nikon Z f camera folder.

Notes
- The profile injects the best current camera→XYZ matrix derived from NX’s PicCon21 into a Nikon Z f DCP. It does not change tone curves or Adobe LUTs.
- For closer matching to NX Studio, you can export an image through Recipe with chroma/hue LUTs applied (scripts/apply_color_matrix.py) and compare.
- As we refine the pipeline (e.g., staged buffer normalization and matrix site), regenerate the matrix (derive_nikon_matrix.py) and rerun `scripts/make_dcp_from_derived.py` to produce an updated profile.
