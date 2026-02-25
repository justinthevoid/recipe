// Package main implements the np3tool binary for VS Code extension IPC.
//
// np3tool reads JSONL messages from stdin and writes JSONL responses to stdout.
// It serves as the backend processing engine for the Recipe NP3 Editor extension.
// stderr is reserved for debug logging (piped to VS Code Output Channel).
//
// CRITICAL: Never write non-JSON to stdout (breaks parser).
// Use stderr for Go's log output.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

// Message represents a JSONL IPC message between the extension host and np3tool.
// All messages follow the shape: { type: "action_name", payload: { ... } }
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// PongPayload is the response payload for np3.ping requests.
type PongPayload struct {
	Status string `json:"status"`
}

// ErrorPayload is the error response payload.
type ErrorPayload struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func main() {
	// Configure structured logging to stderr (debug logging for VS Code Output Channel)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	slog.Info("np3tool started, waiting for JSONL input on stdin")

	scanner := bufio.NewScanner(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for scanner.Scan() {
		line := scanner.Bytes()

		var msg Message
		if err := json.Unmarshal(line, &msg); err != nil {
			slog.Error("Failed to parse input", "error", err)
			sendError(encoder, fmt.Sprintf("malformed JSON: %v", err), "PARSE_ERROR")
			continue
		}

		handleMessage(encoder, &msg)
	}

	if err := scanner.Err(); err != nil {
		slog.Error("Scanner error", "error", err)
	}

	slog.Info("stdin EOF — shutting down gracefully")
}

func handleMessage(encoder *json.Encoder, msg *Message) {
	switch msg.Type {
	case "np3.ping":
		response := Message{
			Type: "np3.pong",
		}
		payload, _ := json.Marshal(PongPayload{Status: "ok"})
		response.Payload = payload

		if err := encoder.Encode(response); err != nil {
			slog.Error("Failed to write pong response", "error", err)
		}

	default:
		slog.Warn("Unknown message type", "type", msg.Type)
		sendError(encoder, fmt.Sprintf("unknown message type: %s", msg.Type), "UNKNOWN_TYPE")
	}
}

func sendError(encoder *json.Encoder, message, code string) {
	payload, _ := json.Marshal(ErrorPayload{
		Message: message,
		Code:    code,
	})

	response := Message{
		Type:    "error",
		Payload: payload,
	}

	if err := encoder.Encode(response); err != nil {
		slog.Error("Failed to write error response", "error", err)
	}
}
