package batch_test

import (
	"runtime"
	"testing"

	"github.com/justin/recipe/internal/batch"
)

func TestConfig_WorkerCount(t *testing.T) {
	tests := []struct {
		name    string
		workers int
		want    int
	}{
		{
			name:    "Default (0)",
			workers: 0,
			want:    runtime.NumCPU(),
		},
		{
			name:    "Auto-detect (-1)",
			workers: -1,
			want:    runtime.NumCPU(),
		},
		{
			name:    "Sequential (1)",
			workers: 1,
			want:    1,
		},
		{
			name:    "Parallel (4)",
			workers: 4,
			want:    4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := batch.Config{
				Workers: tt.workers,
			}
			if got := cfg.WorkerCount(); got != tt.want {
				t.Errorf("WorkerCount() = %v, want %v", got, tt.want)
			}
		})
	}
}
