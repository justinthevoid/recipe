package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/justin/recipe/internal/verify"
	"github.com/spf13/cobra"
	_ "golang.org/x/image/tiff"
)

func newVerifyCmd() *cobra.Command {
	var inputDir string
	var refDir string
	var threshold float64

	const (
		green = "\033[32m"
		red   = "\033[31m"
		reset = "\033[0m"
	)

	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify exported images against a golden reference",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if inputDir == "" {
				return fmt.Errorf("required flag(s) \"input\" not set")
			}
			if refDir == "" {
				return fmt.Errorf("required flag(s) \"reference\" not set")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			comp := verify.NewComparator()
			var failCount int

			// Walk directory
			err := filepath.WalkDir(inputDir, func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}

				// Check extensions
				ext := strings.ToLower(filepath.Ext(path))
				if ext != ".jpg" && ext != ".jpeg" && ext != ".tif" && ext != ".tiff" {
					return nil
				}

				relPath, err := filepath.Rel(inputDir, path)
				if err != nil {
					return err
				}
				refPath := filepath.Join(refDir, relPath)

				// Verify existence
				if _, err := os.Stat(refPath); os.IsNotExist(err) {
					fmt.Printf("%s ... %sFAIL%s (Missing reference)\n", relPath, red, reset)
					failCount++
					return nil
				}

				// Load images
				img1, err := loadImage(path)
				if err != nil {
					fmt.Printf("%s ... %sFAIL%s (Load Input: %v)\n", relPath, red, reset, err)
					failCount++
					return nil
				}
				img2, err := loadImage(refPath)
				if err != nil {
					fmt.Printf("%s ... %sFAIL%s (Load Ref: %v)\n", relPath, red, reset, err)
					failCount++
					return nil
				}

				// Compare
				diff, err := comp.Compare(img1, img2)
				if err != nil {
					fmt.Printf("%s ... %sFAIL%s (Compare: %v)\n", relPath, red, reset, err)
					failCount++
					return nil
				}

				if diff > threshold {
					fmt.Printf("%s ... %sFAIL%s (Diff: %.4f > %.2f)\n", relPath, red, reset, diff, threshold)
					failCount++
				} else {
					fmt.Printf("%s ... %sPASS%s (Diff: %.4f)\n", relPath, green, reset, diff)
				}

				return nil
			})

			if err != nil {
				return err
			}
			if failCount > 0 {
				return fmt.Errorf("%d files failed verification", failCount)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&inputDir, "input", "", "Directory containing exported images to verify")
	cmd.Flags().StringVar(&refDir, "reference", "", "Directory containing golden reference images")
	cmd.Flags().Float64Var(&threshold, "threshold", 0.05, "Comparison threshold (0.0-1.0)")

	_ = cmd.MarkFlagRequired("input")
	_ = cmd.MarkFlagRequired("reference")

	return cmd
}

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}
