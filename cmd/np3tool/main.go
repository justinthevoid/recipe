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
	"io"
	"log/slog"
	"os"

	"github.com/justin/recipe/internal/formats/np3"
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

	// If arguments are provided, run as CLI tool
	if len(os.Args) > 1 {
		runCLI()
		return
	}

	// Otherwise, run as IPC server
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

func runCLI() {
	if os.Args[1] == "get-metadata" {
		var filepath string
		if len(os.Args) > 2 {
			filepath = os.Args[2]
		}

		var data []byte
		var err error

		if filepath == "" || filepath == "-" {
			// Read from stdin
			data, err = io.ReadAll(os.Stdin)
		} else {
			// Read from file
			data, err = os.ReadFile(filepath)
		}

		if err != nil {
			printCLIError("FAILED_READ", fmt.Sprintf("Failed to read input: %v", err))
			os.Exit(1)
		}

		recipe, err := np3.Parse(data)
		if err != nil {
			printCLIError("PARSE_ERROR", fmt.Sprintf("Failed to parse NP3: %v", err))
			os.Exit(1)
		}

		hash := np3.CalculateMagicHash(data)

		// Create metadata response
		metadata := map[string]interface{}{
			"hash":   hash,
			"recipe": recipe,
		}

		response := Message{
			Type:    "np3.metadata",
			Payload: marshalPayload(metadata),
		}

		json.NewEncoder(os.Stdout).Encode(response)
	} else {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func printCLIError(code, message string) {
	errPayload := ErrorPayload{
		Code:    code,
		Message: message,
	}
	response := Message{
		Type:    "error",
		Payload: marshalPayload(errPayload),
	}
	json.NewEncoder(os.Stdout).Encode(response)
}

func marshalPayload(v interface{}) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
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

	case "np3.open":
		var req struct {
			FilePath string `json:"filePath"`
		}
		if err := json.Unmarshal(msg.Payload, &req); err != nil {
			sendError(encoder, "Invalid payload for open", "BAD_REQUEST")
			return
		}

		data, err := os.ReadFile(req.FilePath)
		if err != nil {
			sendError(encoder, fmt.Sprintf("failed to read file: %v", err), "IO_ERROR")
			return
		}

		recipe, err := np3.Parse(data)
		if err != nil {
			switch err {
			case np3.ErrChecksumMismatch:
				sendError(encoder, fmt.Sprintf("failed to parse NP3: %v", err), "ERR_INVALID_CHECKSUM")
			case np3.ErrInvalidMagic:
				sendError(encoder, fmt.Sprintf("failed to parse NP3: %v", err), "ERR_CORRUPTED_FILE")
			default:
				sendError(encoder, fmt.Sprintf("failed to parse NP3: %v", err), "PARSE_ERROR")
			}
			return
		}

		hash := np3.CalculateMagicHash(data)

		metadata := map[string]interface{}{
			"hash":   hash,
			"recipe": recipe,
		}

		response := Message{
			Type:    "np3.metadata",
			Payload: marshalPayload(metadata),
		}

		if err := encoder.Encode(response); err != nil {
			slog.Error("Failed to write metadata response", "error", err)
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
