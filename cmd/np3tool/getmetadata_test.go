package main

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"
)

// TestGetMetadataIPC simulates network/I/O errors (e.g. os.IsPermission) and tests IPC response grace (Task 4).
func TestGetMetadataIPC(t *testing.T) {
	tempDir := t.TempDir()
	badFile := filepath.Join(tempDir, "nonexistent.np3")

	reqMsg := Message{
		Type:    "np3.open",
		Payload: []byte(`{"filePath": "` + badFile + `"}`),
	}

	var out bytes.Buffer
	encoder := json.NewEncoder(&out)

	handleMessage(encoder, &Session{}, &reqMsg)

	var respMsg Message
	if err := json.Unmarshal(out.Bytes(), &respMsg); err != nil {
		t.Fatalf("failed to decode response message: %v", err)
	}

	if respMsg.Type != "error" {
		t.Errorf("expected type 'error', got %q", respMsg.Type)
	}

	var errPayload ErrorPayload
	if err := json.Unmarshal(respMsg.Payload, &errPayload); err != nil {
		t.Fatalf("failed to decode error payload: %v", err)
	}

	if errPayload.Code != "IO_ERROR" {
		t.Errorf("expected error code 'IO_ERROR', got %q", errPayload.Code)
	}
}
