package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestApplyCmd_Run(t *testing.T) {
	// Create test environment
	tmpDir, err := os.MkdirTemp("", "recipe-nx-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")
	np3Path := filepath.Join(tmpDir, "test.np3")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create dummy NP3 file (needs minimum valid structure to pass parser?)
	// The parser checks for magic bytes "NCP" and some length.
	// We'll use a minimal dummy that satisfies np3.ParseMetadata basic check.
	// Magic bytes + version + ... up to minFileSize.
	dummyNP3 := make([]byte, 500)
	copy(dummyNP3, []byte{'N', 'C', 'P'}) // Magic
	copy(dummyNP3[3:], []byte("0200"))    // Version
	if err := os.WriteFile(np3Path, dummyNP3, 0644); err != nil {
		t.Fatal(err)
	}

	// Create dummy NEF
	nefPath := filepath.Join(inputDir, "test.nef")
	if err := os.WriteFile(nefPath, []byte("fake raw data"), 0644); err != nil {
		t.Fatal(err)
	}

	// Run command
	cmd := newApplyCmd()
	// Set output buffer to capture logs? usage?
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	// Invoke via Flags
	cmd.SetArgs([]string{
		"--np3", np3Path,
		"--input", inputDir,
		"--output", outputDir,
		"--manifest", "my-report.json",
	})

	// Execute
	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	// Verify Output
	outNef := filepath.Join(outputDir, "test.nef")
	if _, err := os.Stat(outNef); err != nil {
		t.Errorf("expected output NEF to exist")
	}

	outSidecar := filepath.Join(outputDir, "test.nef.nksc")
	if _, err := os.Stat(outSidecar); err != nil {
		t.Errorf("expected output sidecar to exist")
	}

	// Verify Manifest exists
	manifestPath := filepath.Join(outputDir, "my-report.json")
	if _, err := os.Stat(manifestPath); err != nil {
		t.Errorf("expected my-report.json to exist in output directory")
	}
}

func TestNewApplyCmd(t *testing.T) {
	cmd := newApplyCmd()
	if cmd.Use != "apply" {
		t.Errorf("expected Use to be 'apply', got '%s'", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("expected Short description to be set")
	}

	// Test Flags
	requiredFlags := []string{"np3", "input", "output"}
	for _, flag := range requiredFlags {
		if cmd.Flag(flag) == nil {
			t.Errorf("expected flag --%s to be registered", flag)
		}
	}

	optionalFlags := []string{"overwrite"}
	for _, flag := range optionalFlags {
		if cmd.Flag(flag) == nil {
			t.Errorf("expected flag --%s to be registered", flag)
		}
	}
}

func TestApplyCmd_ExportGuideFlag(t *testing.T) {
	cmd := newApplyCmd()
	if cmd.Flag("export-guide") == nil {
		t.Error("expected flag --export-guide to be registered")
	}
}

func TestApplyCmd_ExportGuideFlow(t *testing.T) {
	// Mock OpenFolder
	orig := openFolderFunc
	defer func() { openFolderFunc = orig }()
	folderOpened := false
	openFolderFunc = func(path string) error {
		folderOpened = true
		return nil
	}

	// Create test environment
	tmpDir, err := os.MkdirTemp("", "recipe-nx-test-guide-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")
	np3Path := filepath.Join(tmpDir, "test.np3")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create dummy NP3 file
	dummyNP3 := make([]byte, 500)
	copy(dummyNP3, []byte{'N', 'C', 'P'})
	copy(dummyNP3[3:], []byte("0200"))
	if err := os.WriteFile(np3Path, dummyNP3, 0644); err != nil {
		t.Fatal(err)
	}

	// Create dummy NEF
	nefPath := filepath.Join(inputDir, "test.nef")
	if err := os.WriteFile(nefPath, []byte("fake raw data"), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := newApplyCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	// Simulate user hitting Enter at the prompt
	cmd.SetIn(bytes.NewBufferString("\n"))

	cmd.SetArgs([]string{
		"--np3", np3Path,
		"--input", inputDir,
		"--output", outputDir,
		"--export-guide",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	output := buf.String()
	if !folderOpened {
		t.Error("expected openFolderFunc to be called")
	}
	if !strings.Contains(output, "Open NX Studio") {
		t.Error("expected export guide instructions containing 'Open NX Studio'")
	}
}

func TestApplyCmd_ExportVerification(t *testing.T) {
	// Mock OpenFolder
	orig := openFolderFunc
	defer func() { openFolderFunc = orig }()
	openFolderFunc = func(path string) error { return nil }

	// Create test environment
	tmpDir, err := os.MkdirTemp("", "recipe-nx-test-verify-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")
	np3Path := filepath.Join(tmpDir, "test.np3")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create dummy NP3 file
	dummyNP3 := make([]byte, 500)
	copy(dummyNP3, []byte{'N', 'C', 'P'})
	copy(dummyNP3[3:], []byte("0200"))
	if err := os.WriteFile(np3Path, dummyNP3, 0644); err != nil {
		t.Fatal(err)
	}

	// Create dummy NEF
	nefPath := filepath.Join(inputDir, "test.nef")
	if err := os.WriteFile(nefPath, []byte("fake raw data"), 0644); err != nil {
		t.Fatal(err)
	}

	// Pre-create the "exported" JPEG
	jpgPath := filepath.Join(outputDir, "test.jpg")
	if err := os.WriteFile(jpgPath, []byte("fake jpg"), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := newApplyCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	// Simulate user hitting Enter
	cmd.SetIn(bytes.NewBufferString("\n"))

	cmd.SetArgs([]string{
		"--np3", np3Path,
		"--input", inputDir,
		"--output", outputDir,
		"--export-guide",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	output := buf.String()
	// Adjust expectation string as we implement it
	expected := "Verification Successful: Found 1/1"
	if !strings.Contains(output, expected) {
		t.Errorf("expected output to contain '%s', got:\n%s", expected, output)
	}
}

func TestApplyCmd_ExportGuideFlow_DryRun(t *testing.T) {
	// Mock OpenFolder - should NOT be called
	orig := openFolderFunc
	defer func() { openFolderFunc = orig }()
	folderOpened := false
	openFolderFunc = func(path string) error {
		folderOpened = true
		return nil
	}

	// Create test environment
	tmpDir, err := os.MkdirTemp("", "recipe-nx-test-dryrun-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")
	np3Path := filepath.Join(tmpDir, "test.np3")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create dummy NP3 file
	dummyNP3 := make([]byte, 500)
	copy(dummyNP3, []byte{'N', 'C', 'P'})
	copy(dummyNP3[3:], []byte("0200"))
	if err := os.WriteFile(np3Path, dummyNP3, 0644); err != nil {
		t.Fatal(err)
	}

	cmd := newApplyCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	cmd.SetArgs([]string{
		"--np3", np3Path,
		"--input", inputDir,
		"--output", outputDir,
		"--export-guide",
		"--dry-run",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	if folderOpened {
		t.Error("expected openFolderFunc NOT to be called in dry-run mode")
	}
	output := buf.String()
	if !strings.Contains(output, "skipped") {
		t.Error("expected message about skipping export guide")
	}
}
