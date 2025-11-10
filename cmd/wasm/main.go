// +build js,wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/justin/recipe/internal/converter"
	"github.com/justin/recipe/internal/formats/lrtemplate"
	"github.com/justin/recipe/internal/formats/np3"
	"github.com/justin/recipe/internal/formats/xmp"
)

// convert exposes the converter.Convert function to JavaScript.
//
// JavaScript signature:
//   convert(inputBytes: Uint8Array, fromFormat: string, toFormat: string) -> Promise<Uint8Array>
//
// Parameters:
//   - inputBytes: Uint8Array containing the source file data
//   - fromFormat: Source format ("np3", "xmp", "lrtemplate", or "" for auto-detect)
//   - toFormat: Target format ("np3", "xmp", "lrtemplate")
//
// Returns:
//   - Promise that resolves to Uint8Array (converted file data)
//   - Promise that rejects with error message string on failure
//
// Example JavaScript usage:
//   const inputData = new Uint8Array(fileBuffer);
//   const outputData = await convert(inputData, "xmp", "np3");
func convertWrapper(this js.Value, args []js.Value) interface{} {
	// Return a Promise to JavaScript
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]

		// Run conversion in a goroutine to avoid blocking
		go func() {
			// Validate arguments (from outer args, not promiseArgs)
			if len(args) < 3 {
				reject.Invoke("convert requires 3 arguments: inputBytes, fromFormat, toFormat")
				return
			}

			// Extract input bytes from Uint8Array (from outer args)
			inputJS := args[0]
			inputLen := inputJS.Get("length").Int()
			inputBytes := make([]byte, inputLen)
			js.CopyBytesToGo(inputBytes, inputJS)

			// Extract format strings (from outer args)
			fromFormat := args[1].String()
			toFormat := args[2].String()

			// Perform conversion
			outputBytes, err := converter.Convert(inputBytes, fromFormat, toFormat)
			if err != nil {
				reject.Invoke(err.Error())
				return
			}

			// Convert output to Uint8Array
			outputJS := js.Global().Get("Uint8Array").New(len(outputBytes))
			js.CopyBytesToJS(outputJS, outputBytes)

			// Resolve promise with output
			resolve.Invoke(outputJS)
		}()

		return nil
	})

	// Create and return Promise
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// detectFormat exposes the converter.DetectFormat function to JavaScript.
//
// JavaScript signature:
//   detectFormat(inputBytes: Uint8Array) -> Promise<string>
//
// Parameters:
//   - inputBytes: Uint8Array containing the file data
//
// Returns:
//   - Promise that resolves to format string ("np3", "xmp", or "lrtemplate")
//   - Promise that rejects with error message string if format cannot be detected
//
// Example JavaScript usage:
//   const inputData = new Uint8Array(fileBuffer);
//   const format = await detectFormat(inputData);
//   console.log("Detected format:", format);
func detectFormatWrapper(this js.Value, args []js.Value) interface{} {
	// Return a Promise to JavaScript
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]

		// Run detection in a goroutine
		go func() {
			// Validate arguments (from outer args, not promiseArgs)
			if len(args) < 1 {
				reject.Invoke("detectFormat requires 1 argument: inputBytes")
				return
			}

			// Extract input bytes from Uint8Array (from outer args)
			inputJS := args[0]
			inputLen := inputJS.Get("length").Int()
			inputBytes := make([]byte, inputLen)
			js.CopyBytesToGo(inputBytes, inputJS)

			// Detect format
			format, err := converter.DetectFormat(inputBytes)
			if err != nil {
				reject.Invoke(err.Error())
				return
			}

			// Resolve promise with format string
			resolve.Invoke(format)
		}()

		return nil
	})

	// Create and return Promise
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// getVersion returns the Recipe version string.
//
// JavaScript signature:
//   getVersion() -> string
//
// Example JavaScript usage:
//   const version = getVersion();
//   console.log("Recipe WASM version:", version);
func getVersionWrapper(this js.Value, args []js.Value) interface{} {
	return "1.0.0-wasm"
}

// extractParameters exposes parameter extraction to JavaScript.
//
// JavaScript signature:
//   extractParameters(inputBytes: Uint8Array, format: string) -> Promise<string>
//
// Parameters:
//   - inputBytes: Uint8Array containing the file data
//   - format: File format ("np3", "xmp", or "lrtemplate")
//
// Returns:
//   - Promise that resolves to JSON string containing extracted parameters
//   - Promise that rejects with error message string if extraction fails
//
// Example JavaScript usage:
//   const inputData = new Uint8Array(fileBuffer);
//   const jsonString = await extractParameters(inputData, "np3");
//   const params = JSON.parse(jsonString);
//   console.log("Exposure:", params.Exposure);
func extractParametersWrapper(this js.Value, args []js.Value) interface{} {
	// Return a Promise to JavaScript
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]

		// Run extraction in a goroutine
		go func() {
			// Validate arguments (from outer args, not promiseArgs)
			if len(args) < 2 {
				reject.Invoke("extractParameters requires 2 arguments: inputBytes, format")
				return
			}

			// Extract input bytes from Uint8Array (from outer args)
			inputJS := args[0]
			inputLen := inputJS.Get("length").Int()
			inputBytes := make([]byte, inputLen)
			js.CopyBytesToGo(inputBytes, inputJS)

			// Extract format string (from outer args)
			format := args[1].String()

			// Parse based on format and extract parameters
			var params map[string]interface{}
			var err error

			switch format {
			case "np3":
				params, err = extractNP3Parameters(inputBytes)
			case "xmp":
				params, err = extractXMPParameters(inputBytes)
			case "lrtemplate":
				params, err = extractLRTemplateParameters(inputBytes)
			default:
				reject.Invoke(fmt.Sprintf("unknown format: %s (expected np3, xmp, or lrtemplate)", format))
				return
			}

			if err != nil {
				reject.Invoke(fmt.Sprintf("parameter extraction failed: %v", err))
				return
			}

			// Convert parameters to JSON string
			jsonBytes, err := json.Marshal(params)
			if err != nil {
				reject.Invoke(fmt.Sprintf("JSON encoding failed: %v", err))
				return
			}

			// Resolve promise with JSON string
			resolve.Invoke(string(jsonBytes))
		}()

		return nil
	})

	// Create and return Promise
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// extractNP3Parameters extracts parameters from NP3 file
func extractNP3Parameters(data []byte) (map[string]interface{}, error) {
	// Parse NP3 file using Epic 1 parser
	recipe, err := np3.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse NP3: %w", err)
	}

	// Convert UniversalRecipe to parameter map
	// Focus on core parameters that users care about
	params := map[string]interface{}{
		"Exposure":    recipe.Exposure,
		"Contrast":    recipe.Contrast,
		"Highlights":  recipe.Highlights,
		"Shadows":     recipe.Shadows,
		"Whites":      recipe.Whites,
		"Blacks":      recipe.Blacks,
		"Vibrance":    recipe.Vibrance,
		"Saturation":  recipe.Saturation,
		"Clarity":     recipe.Clarity,
		"Sharpness":   recipe.Sharpness,
		"Temperature": recipe.Temperature,
		"Tint":        recipe.Tint,
	}

	// Add NP3-specific name if available
	if recipe.Name != "" {
		params["Name"] = recipe.Name
	}

	return params, nil
}

// extractXMPParameters extracts parameters from XMP file
func extractXMPParameters(data []byte) (map[string]interface{}, error) {
	// Parse XMP file using Epic 1 parser
	recipe, err := xmp.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse XMP: %w", err)
	}

	// Convert UniversalRecipe to parameter map
	// Use Lightroom CC naming (Exposure2012, Contrast2012, etc.)
	params := map[string]interface{}{
		"Exposure2012":   recipe.Exposure,
		"Contrast2012":   recipe.Contrast,
		"Highlights2012": recipe.Highlights,
		"Shadows2012":    recipe.Shadows,
		"Whites2012":     recipe.Whites,
		"Blacks2012":     recipe.Blacks,
		"Vibrance":       recipe.Vibrance,
		"Saturation":     recipe.Saturation,
		"Clarity2012":    recipe.Clarity,
		"Sharpness":      recipe.Sharpness,
		"Temperature":    recipe.Temperature,
		"Tint":           recipe.Tint,
		"Dehaze":         recipe.Dehaze,
		"Texture":        recipe.Texture,
	}

	return params, nil
}

// extractLRTemplateParameters extracts parameters from lrtemplate file
func extractLRTemplateParameters(data []byte) (map[string]interface{}, error) {
	// Parse lrtemplate file using Epic 1 parser
	recipe, err := lrtemplate.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse lrtemplate: %w", err)
	}

	// Convert UniversalRecipe to parameter map
	// Use Lightroom Classic naming (same as XMP - both use 2012 process version)
	params := map[string]interface{}{
		"Exposure2012":   recipe.Exposure,
		"Contrast2012":   recipe.Contrast,
		"Highlights2012": recipe.Highlights,
		"Shadows2012":    recipe.Shadows,
		"Whites2012":     recipe.Whites,
		"Blacks2012":     recipe.Blacks,
		"Vibrance":       recipe.Vibrance,
		"Saturation":     recipe.Saturation,
		"Clarity2012":    recipe.Clarity,
		"Sharpness":      recipe.Sharpness,
		"Temperature":    recipe.Temperature,
		"Tint":           recipe.Tint,
		"GrainAmount":    recipe.GrainAmount,
		"GrainSize":      recipe.GrainSize,
	}

	return params, nil
}

func main() {
	// Set up a channel to prevent the program from exiting
	c := make(chan struct{}, 0)

	// Register Go functions as global JavaScript functions
	js.Global().Set("convert", js.FuncOf(convertWrapper))
	js.Global().Set("detectFormat", js.FuncOf(detectFormatWrapper))
	js.Global().Set("extractParameters", js.FuncOf(extractParametersWrapper))
	js.Global().Set("getVersion", js.FuncOf(getVersionWrapper))

	// Signal that WASM is ready
	js.Global().Call("dispatchEvent", js.Global().Get("Event").New("wasmReady"))

	println("Recipe WASM module loaded successfully")
	println("Available functions: convert(), detectFormat(), extractParameters(), getVersion()")

	// Keep the program running
	<-c
}
