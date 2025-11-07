package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegration_BasicConversion(t *testing.T) {
	// Build CLI first
	buildCmd := exec.Command("go", "build", "-o", "recipe-test.exe", ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build CLI: %v", err)
	}
	defer os.Remove("recipe-test.exe")

	// Create temp directory
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.xmp")
	outputFile := filepath.Join(tmpDir, "test.np3")

	// Check if sample file exists
	testDataPath := "../../testdata/xmp/portrait.xmp"
	if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
		t.Skipf("skipping: no test file available at %s", testDataPath)
	}

	// Copy sample file to temp dir
	testData, err := os.ReadFile(testDataPath)
	if err != nil {
		t.Skipf("skipping: cannot read test file (%v)", err)
	}
	os.WriteFile(inputFile, testData, 0644)

	// Run conversion
	cmd := exec.Command("./recipe-test.exe", "convert", inputFile, "--to", "np3")
	output, err := cmd.CombinedOutput()

	// Assertions
	if err != nil {
		t.Fatalf("conversion failed: %v\nOutput: %s", err, output)
	}

	// Check output file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("output file not created")
	}

	// Check success message
	if !strings.Contains(string(output), "✓ Converted") && !strings.Contains(string(output), "Converted") {
		t.Errorf("success message not found in output: %s", output)
	}

	// Check exit code
	if cmd.ProcessState.ExitCode() != 0 {
		t.Errorf("exit code = %d, want 0", cmd.ProcessState.ExitCode())
	}
}

func TestIntegration_FileNotFound(t *testing.T) {
	// Build CLI
	buildCmd := exec.Command("go", "build", "-o", "recipe-test.exe", ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build CLI: %v", err)
	}
	defer os.Remove("recipe-test.exe")

	// Run conversion with nonexistent file
	cmd := exec.Command("./recipe-test.exe", "convert", "nonexistent.xmp", "--to", "np3")
	output, err := cmd.CombinedOutput()

	// Should fail
	if err == nil {
		t.Error("expected command to fail, but it succeeded")
	}

	// Check error message
	if !strings.Contains(string(output), "failed to read input file") {
		t.Errorf("error message not found in output: %s", output)
	}

	// Check exit code
	if cmd.ProcessState.ExitCode() != 1 {
		t.Errorf("exit code = %d, want 1", cmd.ProcessState.ExitCode())
	}
}

func TestIntegration_OverwriteProtection(t *testing.T) {
	buildCmd := exec.Command("go", "build", "-o", "recipe-test.exe", ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build CLI: %v", err)
	}
	defer os.Remove("recipe-test.exe")

	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.xmp")
	outputFile := filepath.Join(tmpDir, "test.np3")

	// Check if sample file exists
	testDataPath := "../../testdata/xmp/portrait.xmp"
	if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
		t.Skipf("skipping: no test file available")
	}

	// Create input file
	testData, _ := os.ReadFile(testDataPath)
	os.WriteFile(inputFile, testData, 0644)

	// Create existing output file
	os.WriteFile(outputFile, []byte("existing"), 0644)

	// Run conversion (should fail)
	cmd := exec.Command("./recipe-test.exe", "convert", inputFile, "--to", "np3")
	output, _ := cmd.CombinedOutput()

	// Check error message
	if !strings.Contains(string(output), "already exists") {
		t.Errorf("overwrite error not found: %s", output)
	}

	// Verify existing file unchanged
	content, _ := os.ReadFile(outputFile)
	if string(content) != "existing" {
		t.Error("existing file was modified")
	}

	// Now test with --overwrite flag
	cmd2 := exec.Command("./recipe-test.exe", "convert", inputFile, "--to", "np3", "--overwrite")
	if err := cmd2.Run(); err != nil {
		t.Fatalf("overwrite failed: %v", err)
	}

	// Verify file was overwritten
	content2, _ := os.ReadFile(outputFile)
	if string(content2) == "existing" {
		t.Error("file was not overwritten despite --overwrite flag")
	}
}

func TestIntegration_MissingToFlag(t *testing.T) {
	buildCmd := exec.Command("go", "build", "-o", "recipe-test.exe", ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build CLI: %v", err)
	}
	defer os.Remove("recipe-test.exe")

	// Run command without --to flag
	cmd := exec.Command("./recipe-test.exe", "convert", "test.xmp")
	output, err := cmd.CombinedOutput()

	// Should fail
	if err == nil {
		t.Error("expected command to fail without --to flag")
	}

	// Check error message mentions required flag
	outputStr := string(output)
	if !strings.Contains(outputStr, "required") || !strings.Contains(outputStr, "to") {
		t.Errorf("required flag error not found in output: %s", output)
	}
}

func TestIntegration_InvalidTargetFormat(t *testing.T) {
	buildCmd := exec.Command("go", "build", "-o", "recipe-test.exe", ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build CLI: %v", err)
	}
	defer os.Remove("recipe-test.exe")

	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.xmp")

	// Create dummy input file
	os.WriteFile(inputFile, []byte("dummy"), 0644)

	// Run command with invalid format
	cmd := exec.Command("./recipe-test.exe", "convert", inputFile, "--to", "pdf")
	output, err := cmd.CombinedOutput()

	// Should fail
	if err == nil {
		t.Error("expected command to fail with invalid format")
	}

	// Check error message
	if !strings.Contains(string(output), "unsupported format") {
		t.Errorf("unsupported format error not found in output: %s", output)
	}
}
