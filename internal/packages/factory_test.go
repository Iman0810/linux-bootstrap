package packages

import (
	"testing"

	"github.com/Iman0810/linux-bootstrap/internal/runner"
)

func TestGetPackageManager(t *testing.T) {
	r := runner.Runner{
		DryRun: true,
	}

	tests := []struct {
		name    string
		manager Manager
		valid   bool
	}{
		{
			name:    "APT",
			manager: APT,
			valid:   true,
		},
		{
			name:    "DNF",
			manager: DNF,
			valid:   true,
		},
		{
			name:    "Pacman",
			manager: Pacman,
			valid:   true,
		},
		{
			name:    "Unknown",
			manager: Unknown,
			valid:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPackageManager(tt.manager, r)

			if tt.valid && got == nil {
				t.Fatalf("expected package manager, got nil")
			}

			if !tt.valid && got != nil {
				t.Fatalf("expected nil, got %T", got)
			}
		})
	}
}