package packages

import "testing"

func TestResolvePackages(t *testing.T) {
	tests := []struct {
		name    string
		manager Manager
		input   []string
		want    []string
	}{
		{
			name:    "apt build tools",
			manager: APT,
			input:   []string{"git", "build-tools", "curl"},
			want:    []string{"git", "build-essential", "curl"},
		},
		{
			name:    "dnf build tools",
			manager: DNF,
			input:   []string{"git", "build-tools", "curl"},
			want:    []string{"git", "gcc gcc-c++ make", "curl"},
		},
		{
			name:    "pacman build tools",
			manager: Pacman,
			input:   []string{"git", "build-tools", "curl"},
			want:    []string{"git", "base-devel", "curl"},
		},
		{
			name:    "unknown manager",
			manager: Unknown,
			input:   []string{"git", "build-tools", "curl"},
			want:    []string{"git", "build-tools", "curl"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolvePackages(tt.manager, tt.input)

			if len(got) != len(tt.want) {
				t.Fatalf(
					"expected %d packages, got %d",
					len(tt.want),
					len(got),
				)
			}

			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Fatalf(
						"package %d: expected %q, got %q",
						i,
						tt.want[i],
						got[i],
					)
				}
			}
		})
	}
}
