package packages

import (
	"testing"

	"github.com/Iman0810/linux-bootstrap/internal/system"
)

func TestDetectManager(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected Manager
	}{
		{
			name:     "Ubuntu",
			id:       "ubuntu",
			expected: APT,
		},
		{
			name:     "Debian",
			id:       "debian",
			expected: APT,
		},
		{
			name:     "Pop OS",
			id:       "pop",
			expected: APT,
		},
		{
			name:     "Fedora",
			id:       "fedora",
			expected: DNF,
		},
		{
			name:     "RHEL",
			id:       "rhel",
			expected: DNF,
		},
		{
			name:     "Arch",
			id:       "arch",
			expected: Pacman,
		},
		{
			name:     "Manjaro",
			id:       "manjaro",
			expected: Pacman,
		},
		{
			name:     "Unknown",
			id:       "something-random",
			expected: Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			osInfo := system.OSInfo{
				ID: tt.id,
			}

			got := DetectManager(osInfo)

			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}