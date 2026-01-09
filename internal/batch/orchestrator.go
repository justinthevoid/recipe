package batch

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/justin/recipe/internal/apperr"
	"github.com/justin/recipe/internal/formats/nksc"
	"github.com/justin/recipe/internal/formats/np3"
)

// Config holds configuration for the batch Orchestrator.
type Config struct {
	InputPath  string
	OutputPath string
	Recipe     *np3.Metadata
	Overwrite  bool
}

// Orchestrator manages the batch processing of files.
type Orchestrator struct {
	Config Config
}

// NewOrchestrator creates a new Orchestrator with the given config.
func NewOrchestrator(cfg Config) *Orchestrator {
	return &Orchestrator{
		Config: cfg,
	}
}

// FindFiles recursively finds all .NEF and .NRW files in the input path.
// Matching is case-insensitive.
// It explicitly excludes files within the OutputPath to avoid processing loops.
func (o *Orchestrator) FindFiles() ([]string, error) {
	var files []string

	absOut, err := filepath.Abs(o.Config.OutputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute output path: %w", err)
	}

	err = filepath.WalkDir(o.Config.InputPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Avoid processing the output directory if it's nested inside input
		absPath, err := filepath.Abs(path)
		if err == nil && strings.HasPrefix(absPath, absOut) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".nef" || ext == ".nrw" {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return files, nil
}

// BatchResult summarizes the outcome of a batch operation.
type BatchResult struct {
	TotalFiles int
	Processed  int
	Skipped    int
	Failed     int
	Errors     []error
}

// ProcessBatch executes the batch processing logic.
func (o *Orchestrator) ProcessBatch() (*BatchResult, error) {
	if o.Config.Recipe == nil {
		return nil, fmt.Errorf("configuration error: NP3 recipe is required for sidecar generation")
	}

	files, err := o.FindFiles()
	if err != nil {
		return nil, err
	}

	logger := slog.Default()

	result := &BatchResult{
		TotalFiles: len(files),
	}

	for _, srcPath := range files {
		fileName := filepath.Base(srcPath)
		fileLogger := logger.With("file", fileName)

		// Calculate relative path to maintain structure
		relPath, err := filepath.Rel(o.Config.InputPath, srcPath)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, apperr.New("resolve_path", fileName, err))
			continue
		}

		dstPath := filepath.Join(o.Config.OutputPath, relPath)
		dstDir := filepath.Dir(dstPath)
		sidecarPath := dstPath + ".nksc"

		// Create output directory structure
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, apperr.New("create_dir", fileName, err))
			continue
		}

		// 1. Handle NEF/Image File copying
		// We ensure the image exists at the destination.
		// Optimization: Only copy if missing or if we want to ensure freshness (not implemented here for speed).
		// Note: The Overwrite flag in Config specifically targets the sidecar generation/overwriting behavior
		// based on the CLI help text ("Overwrite existing NKSC files").
		// However, for robustness, if the NEF is missing, we MUST copy it.
		nefExists := false
		if _, err := os.Stat(dstPath); err == nil {
			nefExists = true
		}

		if !nefExists {
			fileLogger.Debug("copying source image")
			if err := CopyFile(srcPath, dstPath); err != nil {
				result.Failed++
				result.Errors = append(result.Errors, apperr.New("copy_nef", fileName, err))
				continue
			}
		}

		// 2. Handle Sidecar Generation
		// Check if sidecar exists to respect Overwrite flag
		sidecarExists := false
		if _, err := os.Stat(sidecarPath); err == nil {
			sidecarExists = true
		}

		if sidecarExists && !o.Config.Overwrite {
			result.Skipped++
			fileLogger.Debug("skipping existing sidecar")
			continue
		}

		// Generate Sidecar
		fileLogger.Debug("generating sidecar")
		if err := o.generateSidecar(srcPath, dstPath); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, apperr.New("generate_sidecar", fileName, err))
			continue
		}

		result.Processed++
	}

	return result, nil
}

func (o *Orchestrator) generateSidecar(srcPath, targetNEF string) error {
	if o.Config.Recipe == nil {
		return nil
	}

	recipe := nksc.NewNKSCRecipe(o.Config.Recipe, targetNEF)
	xmlBytes, err := recipe.MarshalXML()
	if err != nil {
		return err
	}

	sidecarPath := targetNEF + ".nksc"
	if err := os.WriteFile(sidecarPath, xmlBytes, 0644); err != nil {
		return err
	}
	return nil
}
