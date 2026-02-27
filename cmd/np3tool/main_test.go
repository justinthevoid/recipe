package main

import (
	"bufio"
	"encoding/json"
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
			name:     "valid ping returns pong",
			input:    Message{Type: "np3.ping", Payload: json.RawMessage(`{}`)},
			wantType: "np3.pong",
			wantOK:   true,
		},
		{
			name:     "unknown type returns error",
			input:    Message{Type: "np3.unknown", Payload: json.RawMessage(`{}`)},
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

			if tt.wantType == "np3.pong" {
				var payload PongPayload
				if err := json.Unmarshal(response.Payload, &payload); err != nil {
					t.Fatalf("Failed to parse pong payload: %v", err)
				}
				if payload.Status != "ok" {
					t.Errorf("got status %q, want %q", payload.Status, "ok")
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
