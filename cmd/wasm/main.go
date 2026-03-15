//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/justin/recipe/internal/converter"
	"github.com/justin/recipe/internal/formats/np3"
	"github.com/justin/recipe/internal/formats/xmp"
	"github.com/justin/recipe/internal/lut"
	"github.com/justin/recipe/internal/models"
)

// convert exposes the converter.Convert function to JavaScript.
//
// JavaScript signature:
//
//	convert(inputBytes: Uint8Array, fromFormat: string, toFormat: string) -> Promise<Uint8Array>
//
// Parameters:
//   - inputBytes: Uint8Array containing the source file data
//   - fromFormat: Source format ("np3", "xmp", or "" for auto-detect)
//   - toFormat: Target format ("np3" or "xmp")
func convertWrapper(this js.Value, args []js.Value) interface{} {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]

		go func() {
			if len(args) < 3 {
				reject.Invoke("convert requires 3 arguments: inputBytes, fromFormat, toFormat")
				return
			}

			inputJS := args[0]
			inputLen := inputJS.Get("length").Int()
			inputBytes := make([]byte, inputLen)
			js.CopyBytesToGo(inputBytes, inputJS)

			fromFormat := args[1].String()
			toFormat := args[2].String()

			outputBytes, err := converter.Convert(inputBytes, fromFormat, toFormat)
			if err != nil {
				reject.Invoke(err.Error())
				return
			}

			outputJS := js.Global().Get("Uint8Array").New(len(outputBytes))
			js.CopyBytesToJS(outputJS, outputBytes)
			resolve.Invoke(outputJS)
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// detectFormat exposes the converter.DetectFormat function to JavaScript.
//
// JavaScript signature:
//
//	detectFormat(inputBytes: Uint8Array) -> Promise<string>
func detectFormatWrapper(this js.Value, args []js.Value) interface{} {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]

		go func() {
			if len(args) < 1 {
				reject.Invoke("detectFormat requires 1 argument: inputBytes")
				return
			}

			inputJS := args[0]
			inputLen := inputJS.Get("length").Int()
			inputBytes := make([]byte, inputLen)
			js.CopyBytesToGo(inputBytes, inputJS)

			format, err := converter.DetectFormat(inputBytes)
			if err != nil {
				reject.Invoke(err.Error())
				return
			}

			resolve.Invoke(format)
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// getVersion returns the Recipe version string.
func getVersionWrapper(this js.Value, args []js.Value) interface{} {
	return "1.0.0-wasm"
}

// extractParameters exposes parameter extraction to JavaScript.
//
// JavaScript signature:
//
//	extractParameters(inputBytes: Uint8Array, format: string) -> Promise<string>
func extractParametersWrapper(this js.Value, args []js.Value) interface{} {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]

		go func() {
			if len(args) < 2 {
				reject.Invoke("extractParameters requires 2 arguments: inputBytes, format")
				return
			}

			inputJS := args[0]
			inputLen := inputJS.Get("length").Int()
			inputBytes := make([]byte, inputLen)
			js.CopyBytesToGo(inputBytes, inputJS)

			format := args[1].String()

			var params map[string]interface{}
			var err error

			switch format {
			case "np3":
				params, err = extractNP3Parameters(inputBytes)
			case "xmp":
				params, err = extractXMPParameters(inputBytes)
			default:
				reject.Invoke(fmt.Sprintf("unknown format: %s (expected np3 or xmp)", format))
				return
			}

			if err != nil {
				reject.Invoke(fmt.Sprintf("parameter extraction failed: %v", err))
				return
			}

			jsonBytes, err := json.Marshal(params)
			if err != nil {
				reject.Invoke(fmt.Sprintf("JSON encoding failed: %v", err))
				return
			}

			resolve.Invoke(string(jsonBytes))
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// extractNP3Parameters extracts parameters from NP3 file
func extractNP3Parameters(data []byte) (map[string]interface{}, error) {
	recipe, err := np3.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse NP3: %w", err)
	}

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

	if recipe.Name != "" {
		params["Name"] = recipe.Name
	}

	return params, nil
}

// extractXMPParameters extracts parameters from XMP file
func extractXMPParameters(data []byte) (map[string]interface{}, error) {
	recipe, err := xmp.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse XMP: %w", err)
	}

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

// generate exposes the np3.Generate function to JavaScript.
func generateWrapper(this js.Value, args []js.Value) interface{} {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]

		go func() {
			if len(args) < 1 {
				reject.Invoke("generate requires 1 argument: recipeJSON")
				return
			}

			recipeJSON := args[0].String()

			var recipe models.UniversalRecipe
			if err := json.Unmarshal([]byte(recipeJSON), &recipe); err != nil {
				reject.Invoke(fmt.Sprintf("JSON decode failed: %v", err))
				return
			}

			outputBytes, err := np3.Generate(&recipe)
			if err != nil {
				reject.Invoke(fmt.Sprintf("NP3 generation failed: %v", err))
				return
			}

			outputJS := js.Global().Get("Uint8Array").New(len(outputBytes))
			js.CopyBytesToJS(outputJS, outputBytes)
			resolve.Invoke(outputJS)
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// extractFullRecipe exposes full recipe extraction to JavaScript.
//
// JavaScript signature:
//
//	extractFullRecipe(inputBytes: Uint8Array, format: string) -> Promise<string>
func extractFullRecipeWrapper(this js.Value, args []js.Value) interface{} {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]

		go func() {
			if len(args) < 2 {
				reject.Invoke("extractFullRecipe requires 2 arguments: inputBytes, format")
				return
			}

			inputJS := args[0]
			inputLen := inputJS.Get("length").Int()
			inputBytes := make([]byte, inputLen)
			js.CopyBytesToGo(inputBytes, inputJS)

			format := args[1].String()

			var recipe *models.UniversalRecipe
			var err error

			switch format {
			case "np3":
				recipe, err = np3.Parse(inputBytes)
			case "xmp":
				recipe, err = xmp.Parse(inputBytes)
			default:
				reject.Invoke(fmt.Sprintf("unknown format: %s", format))
				return
			}

			if err != nil {
				reject.Invoke(fmt.Sprintf("parse failed: %v", err))
				return
			}

			jsonBytes, err := json.Marshal(recipe)
			if err != nil {
				reject.Invoke(fmt.Sprintf("JSON encode failed: %v", err))
				return
			}

			resolve.Invoke(string(jsonBytes))
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// generateLUT exposes LUT generation for WebGL preview rendering to JavaScript.
func generateLUTWrapper(this js.Value, args []js.Value) interface{} {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]

		go func() {
			if len(args) < 2 {
				reject.Invoke("generateLUT requires 2 arguments: recipeJSON, size")
				return
			}

			recipeJSON := args[0].String()
			size := args[1].Int()

			var recipe models.UniversalRecipe
			if err := json.Unmarshal([]byte(recipeJSON), &recipe); err != nil {
				reject.Invoke(fmt.Sprintf("JSON decode failed: %v", err))
				return
			}

			lutData, err := lut.Generate3DLUTForPreview(&recipe, size)
			if err != nil {
				reject.Invoke(fmt.Sprintf("LUT generation failed: %v", err))
				return
			}

			outputJS := js.Global().Get("Uint8Array").New(len(lutData))
			js.CopyBytesToJS(outputJS, lutData)

			resolve.Invoke(outputJS)
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("convert", js.FuncOf(convertWrapper))
	js.Global().Set("detectFormat", js.FuncOf(detectFormatWrapper))
	js.Global().Set("extractParameters", js.FuncOf(extractParametersWrapper))
	js.Global().Set("extractFullRecipe", js.FuncOf(extractFullRecipeWrapper))
	js.Global().Set("generate", js.FuncOf(generateWrapper))
	js.Global().Set("generateLUT", js.FuncOf(generateLUTWrapper))
	js.Global().Set("getVersion", js.FuncOf(getVersionWrapper))

	js.Global().Call("dispatchEvent", js.Global().Get("Event").New("wasmReady"))

	println("Recipe WASM module loaded successfully")
	println("Available functions: convert(), detectFormat(), extractParameters(), extractFullRecipe(), generate(), generateLUT(), getVersion()")

	<-c
}
