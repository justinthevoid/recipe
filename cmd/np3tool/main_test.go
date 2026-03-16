package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestPingPong(t *testing.T) {
	binaryPath := buildTestBinary(t)

	tests := []struct {
		name     string
		input    Message
		wantType string
		wantOK   bool
	}{
		{
			name:     "valid ping returns pong with version and echoed requestId",
			input:    Message{Type: "np3.ping", Payload: json.RawMessage(`{}`), RequestId: "req-ping-1"},
			wantType: "np3.pong",
			wantOK:   true,
		},
		{
			name:     "unknown type returns error",
			input:    Message{Type: "np3.unknown", Payload: json.RawMessage(`{}`), RequestId: "req-unknown-1"},
			wantType: "error",
			wantOK:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath)
			cmd.Stderr = io.Discard

			stdin, err := cmd.StdinPipe()
			if err != nil {
				t.Fatalf("Failed to create stdin pipe: %v", err)
			}

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				t.Fatalf("Failed to create stdout pipe: %v", err)
			}

			if err := cmd.Start(); err != nil {
				t.Fatalf("Failed to start binary: %v", err)
			}

			// Send input
			encoder := json.NewEncoder(stdin)
			if err := encoder.Encode(tt.input); err != nil {
				t.Fatalf("Failed to send message: %v", err)
			}
			stdin.Close()

			// Read response
			scanner := bufio.NewScanner(stdout)
			if !scanner.Scan() {
				t.Fatal("No response received from binary")
			}

			var response Message
			if err := json.Unmarshal(scanner.Bytes(), &response); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			if response.Type != tt.wantType {
				t.Errorf("got type %q, want %q", response.Type, tt.wantType)
			}

			// Verify requestId is echoed
			if response.RequestId != tt.input.RequestId {
				t.Errorf("got requestId %q, want %q", response.RequestId, tt.input.RequestId)
			}

			if tt.wantType == "np3.pong" {
				var payload PongPayload
				if err := json.Unmarshal(response.Payload, &payload); err != nil {
					t.Fatalf("Failed to parse pong payload: %v", err)
				}
				if payload.Status != "ok" {
					t.Errorf("got status %q, want %q", payload.Status, "ok")
				}
				if payload.Version == "" {
					t.Error("expected non-empty version field in pong response")
				}
			}

			if tt.wantType == "error" {
				var payload ErrorPayload
				if err := json.Unmarshal(response.Payload, &payload); err != nil {
					t.Fatalf("Failed to parse error payload: %v", err)
				}
				if payload.Code != "UNKNOWN_TYPE" {
					t.Errorf("got code %q, want %q", payload.Code, "UNKNOWN_TYPE")
				}
			}

			if err := cmd.Wait(); err != nil {
				// Expected: process exits after stdin closes
			}
		})
	}
}

func TestMalformedInput(t *testing.T) {
	binaryPath := buildTestBinary(t)

	cmd := exec.Command(binaryPath)
	cmd.Stderr = io.Discard

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Failed to create stdin pipe: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start binary: %v", err)
	}

	// Send malformed JSON
	_, err = stdin.Write([]byte("not valid json\n"))
	if err != nil {
		t.Fatalf("Failed to write to stdin: %v", err)
	}
	stdin.Close()

	// Read error response
	scanner := bufio.NewScanner(stdout)
	if !scanner.Scan() {
		t.Fatal("No response received from binary")
	}

	var response Message
	if err := json.Unmarshal(scanner.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Type != "error" {
		t.Errorf("got type %q, want %q", response.Type, "error")
	}

	var payload ErrorPayload
	if err := json.Unmarshal(response.Payload, &payload); err != nil {
		t.Fatalf("Failed to parse error payload: %v", err)
	}
	if payload.Code != "PARSE_ERROR" {
		t.Errorf("got code %q, want %q", payload.Code, "PARSE_ERROR")
	}

	cmd.Wait()
}

func TestSaveAs(t *testing.T) {
	binaryPath := buildTestBinary(t)

	// Create a temporary file to act as the source
	tmpDir := t.TempDir()
	sourceFile := filepath.Join(tmpDir, "source.np3")
	targetFile := filepath.Join(tmpDir, "target.np3")

	// Standard NP3 file size is 480 bytes
	minimalNP3 := make([]byte, 480)
	copy(minimalNP3, "NCP\x01\x00\x00\x00")
	if err := os.WriteFile(sourceFile, minimalNP3, 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	cmd := exec.Command(binaryPath)
	cmd.Stderr = io.Discard

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Failed to create stdin pipe: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start binary: %v", err)
	}

	encoder := json.NewEncoder(stdin)
	decoder := json.NewDecoder(stdout)

	// 1. Open the file first
	openReq := Message{Type: "np3.open", Payload: marshalPayload(map[string]string{"filePath": sourceFile})}
	if err := encoder.Encode(openReq); err != nil {
		t.Fatalf("Failed to send open message: %v", err)
	}

	var openResp Message
	if err := decoder.Decode(&openResp); err != nil {
		t.Fatalf("Failed to decode open response: %v", err)
	}
	if openResp.Type != "np3.metadata" {
		t.Fatalf("Expected np3.metadata, got %q", openResp.Type)
	}

	// 2. Trigger Save As
	saveAsReq := Message{Type: "np3.save_as", Payload: marshalPayload(map[string]string{"filePath": targetFile})}
	if err := encoder.Encode(saveAsReq); err != nil {
		t.Fatalf("Failed to send save_as message: %v", err)
	}

	var saveAsResp Message
	if err := decoder.Decode(&saveAsResp); err != nil {
		t.Fatalf("Failed to decode save_as response: %v", err)
	}

	if saveAsResp.Type != "np3.save_as_success" {
		t.Fatalf("Expected np3.save_as_success, got %q", saveAsResp.Type)
	}

	// 3. Verify file exists and has content
	if _, err := os.Stat(targetFile); os.IsNotExist(err) {
		t.Errorf("Target file was not created")
	}

	stdin.Close()
	cmd.Wait()
}

func TestGracefulShutdown(t *testing.T) {
	binaryPath := buildTestBinary(t)

	cmd := exec.Command(binaryPath)
	cmd.Stderr = io.Discard

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Failed to create stdin pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start binary: %v", err)
	}

	// Close stdin (EOF) — binary should shut down gracefully
	stdin.Close()

	if err := cmd.Wait(); err != nil {
		t.Errorf("Binary did not exit cleanly after stdin EOF: %v", err)
	}
}

// createTestNP3File creates a minimal valid NP3 test file for integration testing.
func createTestNP3File(t *testing.T) string {
	t.Helper()

	// Use a real NP3 file from testdata
	realFile := filepath.Join("testdata", "Hawthorn.NP3")
	data, err := os.ReadFile(realFile)
	if err != nil {
		t.Fatalf("Failed to read test NP3 file: %v", err)
	}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.np3")
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		t.Fatalf("Failed to write temp NP3 file: %v", err)
	}
	return tmpFile
}

// startBinarySession spawns the binary and opens an NP3 file, returning encoder/decoder.
func startBinarySession(t *testing.T, binaryPath, filePath string) (*exec.Cmd, *json.Encoder, *json.Decoder) {
	t.Helper()

	cmd := exec.Command(binaryPath)
	cmd.Stderr = io.Discard

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Failed to create stdin pipe: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start binary: %v", err)
	}

	encoder := json.NewEncoder(stdin)
	decoder := json.NewDecoder(stdout)

	// Open file
	openReq := Message{Type: "np3.open", Payload: marshalPayload(map[string]string{"filePath": filePath}), RequestId: "open-1"}
	if err := encoder.Encode(openReq); err != nil {
		t.Fatalf("Failed to send open: %v", err)
	}

	var openResp Message
	if err := decoder.Decode(&openResp); err != nil {
		t.Fatalf("Failed to decode open response: %v", err)
	}
	if openResp.Type != "np3.metadata" {
		t.Fatalf("Expected np3.metadata, got %q", openResp.Type)
	}

	return cmd, encoder, decoder
}

func TestPatchValidValue(t *testing.T) {
	binaryPath := buildTestBinary(t)
	testFile := createTestNP3File(t)

	cmd, encoder, decoder := startBinarySession(t, binaryPath, testFile)

	// Send valid patch
	patchReq := Message{
		Type:      "np3.patch",
		RequestId: "patch-1",
		Payload:   marshalPayload(map[string]interface{}{"field": "contrast", "value": 50}),
	}
	if err := encoder.Encode(patchReq); err != nil {
		t.Fatalf("Failed to send patch: %v", err)
	}

	var patchResp Message
	if err := decoder.Decode(&patchResp); err != nil {
		t.Fatalf("Failed to decode patch response: %v", err)
	}

	if patchResp.Type != "np3.patch_success" {
		t.Errorf("Expected np3.patch_success, got %q", patchResp.Type)
	}
	if patchResp.RequestId != "patch-1" {
		t.Errorf("Expected requestId patch-1, got %q", patchResp.RequestId)
	}

	// Check dirty flag
	var payload map[string]interface{}
	json.Unmarshal(patchResp.Payload, &payload)
	if dirty, ok := payload["dirty"].(bool); !ok || !dirty {
		t.Error("Expected dirty: true in patch_success response")
	}

	cmd.Process.Kill()
	cmd.Wait()
}

func TestPatchDoesNotWriteToDisk(t *testing.T) {
	binaryPath := buildTestBinary(t)
	testFile := createTestNP3File(t)

	// Read original content
	originalData, _ := os.ReadFile(testFile)

	cmd, encoder, decoder := startBinarySession(t, binaryPath, testFile)

	// Patch a parameter
	patchReq := Message{
		Type:      "np3.patch",
		RequestId: "patch-nodisk",
		Payload:   marshalPayload(map[string]interface{}{"field": "contrast", "value": 75}),
	}
	encoder.Encode(patchReq)

	var resp Message
	decoder.Decode(&resp)

	if resp.Type != "np3.patch_success" {
		t.Fatalf("Expected np3.patch_success, got %q", resp.Type)
	}

	// Verify file on disk is unchanged
	currentData, _ := os.ReadFile(testFile)
	if string(currentData) != string(originalData) {
		t.Error("np3.patch should NOT write to disk — file was modified")
	}

	cmd.Process.Kill()
	cmd.Wait()
}

func TestPatchOutOfRange(t *testing.T) {
	binaryPath := buildTestBinary(t)
	testFile := createTestNP3File(t)

	cmd, encoder, decoder := startBinarySession(t, binaryPath, testFile)

	// Send out-of-range patch
	patchReq := Message{
		Type:      "np3.patch",
		RequestId: "patch-oor",
		Payload:   marshalPayload(map[string]interface{}{"field": "contrast", "value": 999}),
	}
	encoder.Encode(patchReq)

	var resp Message
	decoder.Decode(&resp)

	if resp.Type != "np3.patch_error" {
		t.Errorf("Expected np3.patch_error, got %q", resp.Type)
	}

	var payload ErrorPayload
	json.Unmarshal(resp.Payload, &payload)
	if payload.Code != "VALIDATION_ERROR" {
		t.Errorf("Expected VALIDATION_ERROR code, got %q", payload.Code)
	}

	cmd.Process.Kill()
	cmd.Wait()
}

func TestPatchNameField(t *testing.T) {
	binaryPath := buildTestBinary(t)
	testFile := createTestNP3File(t)

	cmd, encoder, decoder := startBinarySession(t, binaryPath, testFile)

	// Patch name field with string value — should succeed (F23)
	patchReq := Message{
		Type:      "np3.patch",
		RequestId: "patch-name",
		Payload:   marshalPayload(map[string]interface{}{"field": "name", "value": "My Custom Recipe"}),
	}
	encoder.Encode(patchReq)

	var resp Message
	decoder.Decode(&resp)

	if resp.Type != "np3.patch_success" {
		t.Errorf("Expected np3.patch_success for name field, got %q", resp.Type)
	}

	cmd.Process.Kill()
	cmd.Wait()
}

func TestSaveWritesToDisk(t *testing.T) {
	binaryPath := buildTestBinary(t)
	testFile := createTestNP3File(t)

	originalData, _ := os.ReadFile(testFile)

	cmd, encoder, decoder := startBinarySession(t, binaryPath, testFile)

	// Patch a value
	patchReq := Message{
		Type:      "np3.patch",
		RequestId: "patch-pre-save",
		Payload:   marshalPayload(map[string]interface{}{"field": "contrast", "value": 50}),
	}
	encoder.Encode(patchReq)
	var patchResp Message
	decoder.Decode(&patchResp)

	// Now explicitly save
	saveReq := Message{
		Type:      "np3.save",
		RequestId: "save-1",
		Payload:   json.RawMessage(`{}`),
	}
	encoder.Encode(saveReq)

	var saveResp Message
	decoder.Decode(&saveResp)

	if saveResp.Type != "np3.save_success" {
		t.Errorf("Expected np3.save_success, got %q", saveResp.Type)
	}
	if saveResp.RequestId != "save-1" {
		t.Errorf("Expected requestId save-1, got %q", saveResp.RequestId)
	}

	// Check dirty: false
	var payload map[string]interface{}
	json.Unmarshal(saveResp.Payload, &payload)
	if dirty, ok := payload["dirty"].(bool); !ok || dirty {
		t.Error("Expected dirty: false in save_success response")
	}

	// Verify file on disk changed
	currentData, _ := os.ReadFile(testFile)
	if string(currentData) == string(originalData) {
		t.Error("np3.save should write to disk — file was not modified")
	}

	cmd.Process.Kill()
	cmd.Wait()
}

func TestResetReloadsFromDisk(t *testing.T) {
	binaryPath := buildTestBinary(t)
	testFile := createTestNP3File(t)

	cmd, encoder, decoder := startBinarySession(t, binaryPath, testFile)

	// Patch a value
	patchReq := Message{
		Type:      "np3.patch",
		RequestId: "patch-pre-reset",
		Payload:   marshalPayload(map[string]interface{}{"field": "contrast", "value": 50}),
	}
	encoder.Encode(patchReq)
	var patchResp Message
	decoder.Decode(&patchResp)

	// Reset
	resetReq := Message{
		Type:      "np3.reset",
		RequestId: "reset-1",
		Payload:   json.RawMessage(`{}`),
	}
	encoder.Encode(resetReq)

	var resetResp Message
	decoder.Decode(&resetResp)

	if resetResp.Type != "np3.metadata" {
		t.Errorf("Expected np3.metadata from reset, got %q", resetResp.Type)
	}
	if resetResp.RequestId != "reset-1" {
		t.Errorf("Expected requestId reset-1, got %q", resetResp.RequestId)
	}

	// Check dirty: false
	var payload map[string]interface{}
	json.Unmarshal(resetResp.Payload, &payload)
	if dirty, ok := payload["dirty"].(bool); !ok || dirty {
		t.Error("Expected dirty: false in reset response")
	}

	cmd.Process.Kill()
	cmd.Wait()
}

func TestSyncRequestMemoryOnly(t *testing.T) {
	binaryPath := buildTestBinary(t)
	testFile := createTestNP3File(t)

	originalData, _ := os.ReadFile(testFile)

	cmd, encoder, decoder := startBinarySession(t, binaryPath, testFile)

	// Get the recipe from the open response metadata
	// Send sync_request
	syncReq := Message{
		Type:      "np3.sync_request",
		RequestId: "sync-1",
		Payload:   marshalPayload(map[string]interface{}{"recipe": map[string]interface{}{"contrast": 99}}),
	}
	encoder.Encode(syncReq)

	var syncResp Message
	decoder.Decode(&syncResp)

	if syncResp.Type != "np3.sync" {
		t.Errorf("Expected np3.sync, got %q", syncResp.Type)
	}

	// Verify file on disk is unchanged (F11)
	currentData, _ := os.ReadFile(testFile)
	if string(currentData) != string(originalData) {
		t.Error("np3.sync_request should NOT write to disk — file was modified")
	}

	cmd.Process.Kill()
	cmd.Wait()
}

func TestRapidPatchRequestIds(t *testing.T) {
	binaryPath := buildTestBinary(t)
	testFile := createTestNP3File(t)

	cmd, encoder, decoder := startBinarySession(t, binaryPath, testFile)

	// Send 10 rapid patches
	for i := 0; i < 10; i++ {
		patchReq := Message{
			Type:      "np3.patch",
			RequestId: fmt.Sprintf("rapid-%d", i),
			Payload:   marshalPayload(map[string]interface{}{"field": "contrast", "value": float64(i)}),
		}
		encoder.Encode(patchReq)
	}

	// Read 10 responses and verify each has correct requestId
	receivedIds := make(map[string]bool)
	for i := 0; i < 10; i++ {
		var resp Message
		if err := decoder.Decode(&resp); err != nil {
			t.Fatalf("Failed to decode response %d: %v", i, err)
		}
		if resp.Type != "np3.patch_success" {
			t.Errorf("Response %d: expected np3.patch_success, got %q", i, resp.Type)
		}
		receivedIds[resp.RequestId] = true
	}

	// Verify all 10 requestIds were received
	for i := 0; i < 10; i++ {
		expectedId := fmt.Sprintf("rapid-%d", i)
		if !receivedIds[expectedId] {
			t.Errorf("Missing response for requestId %q", expectedId)
		}
	}

	cmd.Process.Kill()
	cmd.Wait()
}

func TestSaveAsClearsDirty(t *testing.T) {
	binaryPath := buildTestBinary(t)
	testFile := createTestNP3File(t)

	cmd, encoder, decoder := startBinarySession(t, binaryPath, testFile)

	// Patch to make dirty
	patchReq := Message{
		Type:      "np3.patch",
		RequestId: "patch-before-saveas",
		Payload:   marshalPayload(map[string]interface{}{"field": "contrast", "value": 25}),
	}
	encoder.Encode(patchReq)
	var patchResp Message
	decoder.Decode(&patchResp)

	// Save As
	targetFile := filepath.Join(t.TempDir(), "saveas-target.np3")
	saveAsReq := Message{
		Type:      "np3.save_as",
		RequestId: "saveas-dirty",
		Payload:   marshalPayload(map[string]string{"filePath": targetFile}),
	}
	encoder.Encode(saveAsReq)

	var saveAsResp Message
	decoder.Decode(&saveAsResp)

	if saveAsResp.Type != "np3.save_as_success" {
		t.Fatalf("Expected np3.save_as_success, got %q", saveAsResp.Type)
	}

	// Check dirty: false (F12)
	var payload map[string]interface{}
	json.Unmarshal(saveAsResp.Payload, &payload)
	if dirty, ok := payload["dirty"].(bool); !ok || dirty {
		t.Error("Expected dirty: false in save_as_success response")
	}

	cmd.Process.Kill()
	cmd.Wait()
}

func TestCLIGetMetadata(t *testing.T) {
	binaryPath := buildTestBinary(t)
	testFile := createTestNP3File(t)

	// Run CLI mode: get-metadata with a real NP3 file
	cmd := exec.Command(binaryPath, "get-metadata", testFile)
	cmd.Stderr = io.Discard

	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("CLI get-metadata failed: %v", err)
	}

	var response Message
	if err := json.Unmarshal(output, &response); err != nil {
		t.Fatalf("Failed to parse CLI output: %v", err)
	}

	if response.Type != "np3.metadata" {
		t.Errorf("Expected type np3.metadata, got %q", response.Type)
	}

	// Verify payload contains recipe and parameterDefinitions
	var payload map[string]interface{}
	if err := json.Unmarshal(response.Payload, &payload); err != nil {
		t.Fatalf("Failed to parse payload: %v", err)
	}

	if _, ok := payload["recipe"]; !ok {
		t.Error("Missing 'recipe' in CLI get-metadata output")
	}
	if _, ok := payload["parameterDefinitions"]; !ok {
		t.Error("Missing 'parameterDefinitions' in CLI get-metadata output")
	}
	if _, ok := payload["hash"]; !ok {
		t.Error("Missing 'hash' in CLI get-metadata output")
	}
}

func TestCLIUnknownCommand(t *testing.T) {
	binaryPath := buildTestBinary(t)

	cmd := exec.Command(binaryPath, "nonexistent-command")
	err := cmd.Run()
	if err == nil {
		t.Error("Expected non-zero exit code for unknown CLI command")
	}
}

func TestFilePermissionsPreserved(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("File permissions not applicable on Windows")
	}

	binaryPath := buildTestBinary(t)
	testFile := createTestNP3File(t)

	// Set custom permissions
	if err := os.Chmod(testFile, 0600); err != nil {
		t.Fatalf("Failed to set permissions: %v", err)
	}

	cmd, encoder, decoder := startBinarySession(t, binaryPath, testFile)

	// Patch and save
	patchReq := Message{
		Type:      "np3.patch",
		RequestId: "patch-perm",
		Payload:   marshalPayload(map[string]interface{}{"field": "contrast", "value": 25}),
	}
	encoder.Encode(patchReq)
	var patchResp Message
	decoder.Decode(&patchResp)

	saveReq := Message{
		Type:      "np3.save",
		RequestId: "save-perm",
		Payload:   json.RawMessage(`{}`),
	}
	encoder.Encode(saveReq)
	var saveResp Message
	decoder.Decode(&saveResp)

	if saveResp.Type != "np3.save_success" {
		t.Fatalf("Expected np3.save_success, got %q", saveResp.Type)
	}

	// Verify permissions preserved (P2-10b)
	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("Expected permissions 0600, got %o", info.Mode().Perm())
	}

	cmd.Process.Kill()
	cmd.Wait()
}

func buildTestBinary(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	binaryName := "np3tool"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binaryPath := filepath.Join(tmpDir, binaryName)

	// Build from the source directory
	srcDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Dir = srcDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build test binary: %v\n%s", err, output)
	}

	return binaryPath
}
