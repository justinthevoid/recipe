package main

import (
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
)

// filesLoadedMsg is sent when files are loaded
type filesLoadedMsg struct {
	files []FileInfo
	err   error
}

// loadFilesCmd loads files from a directory
func loadFilesCmd(dir string) tea.Cmd {
	return func() tea.Msg {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return filesLoadedMsg{err: err}
		}

		var files []FileInfo
		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			// Filter: only directories and preset files
			ext := filepath.Ext(entry.Name())
			isPreset := ext == ".np3" || ext == ".xmp" || ext == ".lrtemplate" || ext == ".costyle" || ext == ".costylepack"

			if entry.IsDir() || isPreset {
				files = append(files, FileInfo{
					Name:    entry.Name(),
					Path:    filepath.Join(dir, entry.Name()),
					Size:    info.Size(),
					ModTime: info.ModTime(),
					IsDir:   entry.IsDir(),
					Format:  detectFormat(entry.Name(), entry.IsDir()),
				})
			}
		}

		return filesLoadedMsg{files: files}
	}
}

// detectFormat determines the file format based on extension
func detectFormat(name string, isDir bool) string {
	if isDir {
		return "dir"
	}

	ext := filepath.Ext(name)
	switch ext {
	case ".np3":
		return "np3"
	case ".xmp":
		return "xmp"
	case ".lrtemplate":
		return "lrtemplate"
	case ".costyle":
		return "costyle"
	case ".costylepack":
		return "costylepack"
	default:
		return ""
	}
}
