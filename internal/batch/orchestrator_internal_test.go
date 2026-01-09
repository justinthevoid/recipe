package batch

import (
	"testing"
	"time"
)

func TestOrchestrator_shouldSkip(t *testing.T) {
	startTime := time.Date(2026, 1, 8, 10, 0, 0, 0, time.UTC)

	orch := &Orchestrator{
		Config: Config{Force: false},
		previousResults: map[string]FileResult{
			"test.nef": {
				InputPath:   "test.nef",
				Status:      "success",
				NP3Hash:     "hash_v1",
				Size:        1000,
				ModTime:     startTime,
				PayloadHash: "nef_hash_1",
			},
			"failed.nef": {
				InputPath: "failed.nef",
				Status:    "error",
			},
		},
	}

	tests := []struct {
		name        string
		inputPath   string
		currentNP3  string
		currentSize int64
		currentMod  time.Time
		force       bool
		expectSkip  bool
		desc        string
	}{
		{
			name:        "Exact Match",
			inputPath:   "test.nef",
			currentNP3:  "hash_v1",
			currentSize: 1000,
			currentMod:  startTime,
			force:       false,
			expectSkip:  true,
			desc:        "Should skip when all attributes match",
		},
		{
			name:        "NP3 Changed",
			inputPath:   "test.nef",
			currentNP3:  "hash_v2", // Different
			currentSize: 1000,
			currentMod:  startTime,
			force:       false,
			expectSkip:  false,
			desc:        "Should NOT skip when NP3 hash changes",
		},
		{
			name:        "File Size Changed",
			inputPath:   "test.nef",
			currentNP3:  "hash_v1",
			currentSize: 1001, // Different
			currentMod:  startTime,
			force:       false,
			expectSkip:  false,
			desc:        "Should NOT skip when file size changes",
		},
		{
			name:        "ModTime Changed",
			inputPath:   "test.nef",
			currentNP3:  "hash_v1",
			currentSize: 1000,
			currentMod:  startTime.Add(1 * time.Second), // Different
			force:       false,
			expectSkip:  false,
			desc:        "Should NOT skip when ModTime changes",
		},
		{
			name:        "Force Enabled",
			inputPath:   "test.nef",
			currentNP3:  "hash_v1",
			currentSize: 1000,
			currentMod:  startTime,
			force:       true, // FORCE
			expectSkip:  false,
			desc:        "Should NOT skip when Force is true",
		},
		{
			name:        "Previously Failed",
			inputPath:   "failed.nef",
			currentNP3:  "hash_v1",
			currentSize: 1000,
			currentMod:  startTime,
			force:       false,
			expectSkip:  false,
			desc:        "Should NOT skip if previous status was error",
		},
		{
			name:        "New File",
			inputPath:   "new.nef",
			currentNP3:  "hash_v1",
			currentSize: 1000,
			currentMod:  startTime,
			force:       false,
			expectSkip:  false,
			desc:        "Should NOT skip if file not in manifest",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			orch.Config.Force = tc.force
			skip := orch.shouldSkip(tc.inputPath, tc.currentNP3, tc.currentSize, tc.currentMod)
			if skip != tc.expectSkip {
				t.Errorf("%s: expected skip=%v, got %v", tc.desc, tc.expectSkip, skip)
			}
		})
	}
}
