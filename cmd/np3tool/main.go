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
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"strings"

	"github.com/justin/recipe/internal/formats/np3"
	"github.com/justin/recipe/internal/models"
)

// Message represents a JSONL IPC message between the extension host and np3tool.
// All messages follow the shape: { type: "action_name", payload: { ... } }
type Message struct {
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	RequestId string          `json:"requestId,omitempty"`
}

// PongPayload is the response payload for np3.ping requests.
type PongPayload struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// ErrorPayload is the error response payload.
type ErrorPayload struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	RawData string `json:"rawData,omitempty"`
}

// Session holds state for the currently open file
type Session struct {
	FilePath string
	Recipe   *models.UniversalRecipe
	Dirty    bool
}

// ValidationError is returned by field validation.
type ValidationError struct {
	Code    string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// Error code registry with user-facing messages and suggested actions (P2-7a).
var errorMessages = map[string]struct {
	UserMessage string
	Action      string
}{
	"PARSE_ERROR":         {UserMessage: "The file could not be parsed.", Action: "Ensure the file is a valid NP3 file."},
	"BAD_REQUEST":         {UserMessage: "Invalid request format.", Action: "This is a bug — please report it."},
	"BAD_STATE":           {UserMessage: "No file is currently open.", Action: "Open an NP3 file first."},
	"IO_ERROR":            {UserMessage: "File read/write failed.", Action: "Check file permissions and disk space."},
	"GENERATE_ERROR":      {UserMessage: "Failed to generate NP3 binary.", Action: "The recipe may contain invalid data."},
	"PATCH_ERROR":         {UserMessage: "Failed to apply the edit.", Action: "Try undoing and re-applying."},
	"INVALID_PATH":        {UserMessage: "The file path is empty or invalid.", Action: "Choose a valid file location."},
	"UNKNOWN_TYPE":        {UserMessage: "Unrecognized message type.", Action: "Extension may need updating."},
	"ERR_INVALID_CHECKSUM": {UserMessage: "NP3 file checksum validation failed.", Action: "The file may be corrupted. Try the .bak backup."},
	"ERR_CORRUPTED_FILE":  {UserMessage: "NP3 magic bytes are invalid.", Action: "The file is not a valid NP3 file."},
	"VALIDATION_ERROR":    {UserMessage: "Parameter value is out of range.", Action: "Use a value within the allowed range."},
	"UNKNOWN_FIELD":       {UserMessage: "Unknown parameter field.", Action: "This field is not recognized by NP3 format."},
}


// validatePatchField validates that a field exists and its value is in range.
// Metadata fields (name, description, sourceFormat) are exempt from range checks.
func validatePatchField(field string, value any) error {
	// Metadata fields are always valid (F23)
	switch field {
	case "name", "description", "sourceFormat":
		return nil
	}

	defs := models.GetNP3ParameterDefinitions()
	for _, def := range defs {
		if def.Key == field {
			// Must be numeric
			numVal, ok := toFloat64(value)
			if !ok {
				return &ValidationError{Code: "VALIDATION_ERROR", Message: fmt.Sprintf("field %q requires a numeric value", field)}
			}
			if numVal < def.Min || numVal > def.Max {
				return &ValidationError{Code: "VALIDATION_ERROR", Message: fmt.Sprintf("field %q value %v is out of range [%v, %v]", field, numVal, def.Min, def.Max)}
			}
			return nil
		}
	}

	return &ValidationError{Code: "UNKNOWN_FIELD", Message: fmt.Sprintf("unknown field: %q", field)}
}

func toFloat64(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	default:
		return 0, false
	}
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
	session := &Session{}

	for scanner.Scan() {
		line := scanner.Bytes()

		var msg Message
		if err := json.Unmarshal(line, &msg); err != nil {
			slog.Error("Failed to parse input", "error", err)
			sendErrorWithId(encoder, "", "error", fmt.Sprintf("malformed JSON: %v", err), "PARSE_ERROR")
			continue
		}

		handleMessage(encoder, session, &msg)
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
			"hash":                 hash,
			"recipe":               recipe,
			"parameterDefinitions": models.GetNP3ParameterDefinitions(),
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
	b, err := json.Marshal(v)
	if err != nil {
		slog.Error("Failed to marshal payload", "error", err)
		return json.RawMessage(`{"error":"marshal_failed"}`)
	}
	return b
}

func handleMessage(encoder *json.Encoder, session *Session, msg *Message) {
	switch msg.Type {
	case "np3.ping":
		response := Message{
			Type:      "np3.pong",
			RequestId: msg.RequestId,
		}
		payload, _ := json.Marshal(PongPayload{Status: "ok", Version: "1.0.0"})
		response.Payload = payload

		if err := encoder.Encode(response); err != nil {
			slog.Error("Failed to write pong response", "error", err)
		}

	case "np3.open":
		var req struct {
			FilePath string `json:"filePath"`
		}
		if err := json.Unmarshal(msg.Payload, &req); err != nil {
			sendErrorWithId(encoder, msg.RequestId, "error", "Invalid payload for open", "BAD_REQUEST")
			return
		}

		data, err := os.ReadFile(req.FilePath)
		if err != nil {
			sendErrorWithId(encoder, msg.RequestId, "error", fmt.Sprintf("failed to read file: %v", err), "IO_ERROR")
			return
		}

		recipe, err := np3.Parse(data)
		if err != nil {
			encodedData := base64.StdEncoding.EncodeToString(data)
			switch err {
			case np3.ErrChecksumMismatch:
				sendErrorWithDataAndId(encoder, msg.RequestId, fmt.Sprintf("failed to parse NP3: %v", err), "ERR_INVALID_CHECKSUM", encodedData)
			case np3.ErrInvalidMagic:
				sendErrorWithDataAndId(encoder, msg.RequestId, fmt.Sprintf("failed to parse NP3: %v", err), "ERR_CORRUPTED_FILE", encodedData)
			default:
				sendErrorWithDataAndId(encoder, msg.RequestId, fmt.Sprintf("failed to parse NP3: %v", err), "PARSE_ERROR", encodedData)
			}
			return
		}

		session.FilePath = req.FilePath
		session.Recipe = recipe
		session.Dirty = false

		hash := np3.CalculateMagicHash(data)

		metadata := map[string]interface{}{
			"hash":                 hash,
			"recipe":               recipe,
			"parameterDefinitions": models.GetNP3ParameterDefinitions(),
			"dirty":                false,
		}

		response := Message{
			Type:      "np3.metadata",
			Payload:   marshalPayload(metadata),
			RequestId: msg.RequestId,
		}

		if err := encoder.Encode(response); err != nil {
			slog.Error("Failed to write metadata response", "error", err)
		}

	case "np3.patch":
		if session.FilePath == "" || session.Recipe == nil {
			sendErrorWithId(encoder, msg.RequestId, "np3.patch_error", "No file is currently open", "BAD_STATE")
			return
		}

		var req struct {
			Field string `json:"field"`
			Value any    `json:"value"`
		}
		if err := json.Unmarshal(msg.Payload, &req); err != nil {
			sendErrorWithId(encoder, msg.RequestId, "np3.patch_error", "Invalid payload for patch", "BAD_REQUEST")
			return
		}

		// Validate parameter (P1-2b)
		if err := validatePatchField(req.Field, req.Value); err != nil {
			var ve *ValidationError
			code := "VALIDATION_ERROR"
			if errors.As(err, &ve) {
				code = ve.Code
			}
			sendErrorWithId(encoder, msg.RequestId, "np3.patch_error", err.Error(), code)
			return
		}

		// Dynamically update field via JSON round-trip
		b, _ := json.Marshal(session.Recipe)
		var m map[string]interface{}
		json.Unmarshal(b, &m)

		// Support nested paths like "red.hue"
		parts := strings.Split(req.Field, ".")
		curr := m
		for i := 0; i < len(parts)-1; i++ {
			if next, ok := curr[parts[i]].(map[string]interface{}); ok {
				curr = next
			} else {
				newMap := make(map[string]interface{})
				curr[parts[i]] = newMap
				curr = newMap
			}
		}
		curr[parts[len(parts)-1]] = req.Value

		b2, _ := json.Marshal(m)
		if err := json.Unmarshal(b2, session.Recipe); err != nil {
			sendErrorWithId(encoder, msg.RequestId, "np3.patch_error", fmt.Sprintf("failed to apply patch: %v", err), "PATCH_ERROR")
			return
		}

		// Memory-only — no disk write (P1-2a)
		session.Dirty = true

		response := Message{
			Type:      "np3.patch_success",
			RequestId: msg.RequestId,
			Payload: marshalPayload(map[string]interface{}{
				"field": req.Field,
				"value": req.Value,
				"dirty": true,
			}),
		}
		if err := encoder.Encode(response); err != nil {
			slog.Error("Failed to write patch_success response", "error", err)
		}

	case "np3.save":
		if session.FilePath == "" || session.Recipe == nil {
			sendErrorWithId(encoder, msg.RequestId, "np3.save_error", "No file is currently open", "BAD_STATE")
			return
		}

		if err := saveSession(session); err != nil {
			sendErrorWithId(encoder, msg.RequestId, "np3.save_error", err.Error(), "IO_ERROR")
			return
		}

		session.Dirty = false

		response := Message{
			Type:      "np3.save_success",
			RequestId: msg.RequestId,
			Payload: marshalPayload(map[string]interface{}{
				"filePath": session.FilePath,
				"dirty":    false,
			}),
		}
		if err := encoder.Encode(response); err != nil {
			slog.Error("Failed to write save_success response", "error", err)
		}

	case "np3.save_as":
		var req struct {
			FilePath string `json:"filePath"`
		}
		if err := json.Unmarshal(msg.Payload, &req); err != nil {
			sendErrorWithId(encoder, msg.RequestId, "error", "Invalid payload for save_as", "BAD_REQUEST")
			return
		}

		if req.FilePath == "" {
			sendErrorWithId(encoder, msg.RequestId, "error", "Save path cannot be empty", "INVALID_PATH")
			return
		}

		if session.Recipe == nil {
			sendErrorWithId(encoder, msg.RequestId, "error", "No active recipe session to save", "BAD_STATE")
			return
		}

		data, err := np3.Generate(session.Recipe)
		if err != nil {
			sendErrorWithId(encoder, msg.RequestId, "error", fmt.Sprintf("failed to generate NP3: %v", err), "GENERATE_ERROR")
			return
		}

		// Preserve source file permissions for save_as, fallback to 0644 (P2-10b)
		perm := os.FileMode(0644)
		if session.FilePath != "" {
			if info, err := os.Stat(session.FilePath); err == nil {
				perm = info.Mode().Perm()
			}
		}
		if err := os.WriteFile(req.FilePath, data, perm); err != nil {
			sendErrorWithId(encoder, msg.RequestId, "error", fmt.Sprintf("failed to write file: %v", err), "IO_ERROR")
			return
		}

		// Update session to the new file path and clear dirty (F12)
		session.FilePath = req.FilePath
		session.Dirty = false

		response := Message{
			Type:      "np3.save_as_success",
			RequestId: msg.RequestId,
			Payload: marshalPayload(map[string]interface{}{
				"filePath": req.FilePath,
				"dirty":    false,
			}),
		}
		if err := encoder.Encode(response); err != nil {
			slog.Error("Failed to write save_as_success response", "error", err)
		}

	case "np3.sync_request":
		if session.FilePath == "" {
			sendErrorWithId(encoder, msg.RequestId, "np3.sync_error", "No file is currently open", "BAD_STATE")
			return
		}

		var req struct {
			Recipe *models.UniversalRecipe `json:"recipe"`
		}
		if err := json.Unmarshal(msg.Payload, &req); err != nil {
			sendErrorWithId(encoder, msg.RequestId, "np3.sync_error", "Invalid payload for sync_request", "BAD_REQUEST")
			return
		}

		// Memory-only — no disk write (F11)
		session.Recipe = req.Recipe
		session.Dirty = true

		// Echo full recipe payload with dirty flag (F34)
		response := Message{
			Type:      "np3.sync",
			RequestId: msg.RequestId,
			Payload: marshalPayload(map[string]interface{}{
				"recipe": session.Recipe,
				"dirty":  true,
			}),
		}
		if err := encoder.Encode(response); err != nil {
			slog.Error("Failed to write sync response", "error", err)
		}

	case "np3.reset":
		if session.FilePath == "" || session.Recipe == nil {
			sendErrorWithId(encoder, msg.RequestId, "error", "No file is currently open", "BAD_STATE")
			return
		}

		data, err := os.ReadFile(session.FilePath)
		if err != nil {
			sendErrorWithId(encoder, msg.RequestId, "error", "Original file not found. It may have been moved or deleted.", "IO_ERROR")
			return
		}

		recipe, err := np3.Parse(data)
		if err != nil {
			sendErrorWithId(encoder, msg.RequestId, "error", fmt.Sprintf("failed to parse NP3 on reset: %v", err), "PARSE_ERROR")
			return
		}

		session.Recipe = recipe
		session.Dirty = false

		hash := np3.CalculateMagicHash(data)

		metadata := map[string]interface{}{
			"hash":                 hash,
			"recipe":               recipe,
			"parameterDefinitions": models.GetNP3ParameterDefinitions(),
			"dirty":                false,
		}

		response := Message{
			Type:      "np3.metadata",
			RequestId: msg.RequestId,
			Payload:   marshalPayload(metadata),
		}

		if err := encoder.Encode(response); err != nil {
			slog.Error("Failed to write reset response", "error", err)
		}

	default:
		slog.Warn("Unknown message type", "type", msg.Type)
		sendErrorWithId(encoder, msg.RequestId, "error", fmt.Sprintf("unknown message type: %s", msg.Type), "UNKNOWN_TYPE")
	}
}

func sendErrorWithId(encoder *json.Encoder, requestId, msgType, message, code string) {
	payload, _ := json.Marshal(ErrorPayload{
		Message: message,
		Code:    code,
	})

	response := Message{
		Type:      msgType,
		Payload:   payload,
		RequestId: requestId,
	}

	if err := encoder.Encode(response); err != nil {
		slog.Error("Failed to write error response", "error", err)
	}
}

func sendErrorWithDataAndId(encoder *json.Encoder, requestId, message, code, rawData string) {
	payload, _ := json.Marshal(ErrorPayload{
		Message: message,
		Code:    code,
		RawData: rawData,
	})

	response := Message{
		Type:      "error",
		Payload:   payload,
		RequestId: requestId,
	}

	if err := encoder.Encode(response); err != nil {
		slog.Error("Failed to write error response", "error", err)
	}
}

// saveSession writes the current session recipe to disk,
// preserving original file permissions (P2-10b).
func saveSession(session *Session) error {
	data, err := np3.Generate(session.Recipe)
	if err != nil {
		return fmt.Errorf("failed to generate NP3: %w", err)
	}

	// Preserve original file permissions
	perm := os.FileMode(0644)
	if info, err := os.Stat(session.FilePath); err == nil {
		perm = info.Mode().Perm()
	}

	if err := os.WriteFile(session.FilePath, data, perm); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
