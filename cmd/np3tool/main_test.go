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
