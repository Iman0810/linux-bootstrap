package profile

import (
	"testing"

	"github.com/Iman0810/linux-bootstrap/internal/packages"
)

func TestGetProfile(t *testing.T) {
	tests := []struct {
		name       string
		profile    string
		shouldFind bool
	}{
		{
			name:       "Essentials",
			profile:    "essentials",
			shouldFind: true,
		},
		{
			name:       "Development",
			profile:    "development",
			shouldFind: true,
		},
		{
			name:       "Multimedia",
			profile:    "multimedia",
			shouldFind: true,
		},
		{
			name:       "Unknown profile",
			profile:    "banana",
			shouldFind: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, ok := Get(tt.profile)

			if ok != tt.shouldFind {
				t.Fatalf(
					"expected found=%v, got %v",
					tt.shouldFind,
					ok,
				)
			}

			if tt.shouldFind && profile.Name != tt.profile {
				t.Fatalf(
					"expected profile name %q, got %q",
					tt.profile,
					profile.Name,
				)
			}
		})
	}
}

func TestProfilePackages(t *testing.T) {
	tests := []struct {
		name     string
		profile  Profile
		manager  packages.Manager
		expected []string
	}{
		{
			name:    "Essentials APT",
			profile: Essentials,
			manager: packages.APT,
			expected: []string{
				"git",
				"curl",
				"wget",
				"unzip",
			},
		},
		{
			name:    "Essentials DNF",
			profile: Essentials,
			manager: packages.DNF,
			expected: []string{
				"git",
				"curl",
				"wget",
				"unzip",
			},
		},
		{
			name:    "Essentials Pacman",
			profile: Essentials,
			manager: packages.Pacman,
			expected: []string{
				"git",
				"curl",
				"wget",
				"unzip",
			},
		},
		{
			name:    "Development APT",
			profile: Development,
			manager: packages.APT,
			expected: []string{
				"git",
				"curl",
				"wget",
				"unzip",
				"build-essential",
			},
		},
		{
			name:    "Development DNF",
			profile: Development,
			manager: packages.DNF,
			expected: []string{
				"git",
				"curl",
				"wget",
				"unzip",
				"gcc",
				"gcc-c++",
				"make",
			},
		},
		{
			name:    "Development Pacman",
			profile: Development,
			manager: packages.Pacman,
			expected: []string{
				"git",
				"curl",
				"wget",
				"unzip",
				"base-devel",
			},
		},
		{
			name:    "Multimedia APT",
			profile: Multimedia,
			manager: packages.APT,
			expected: []string{
				"ffmpeg",
				"vlc",
			},
		},
		{
			name:    "Multimedia DNF",
			profile: Multimedia,
			manager: packages.DNF,
			expected: []string{
				"ffmpeg",
				"vlc",
			},
		},
		{
			name:    "Multimedia Pacman",
			profile: Multimedia,
			manager: packages.Pacman,
			expected: []string{
				"ffmpeg",
				"vlc",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, ok := PackagesFor(tt.profile, tt.manager)

			if !ok {
				t.Fatalf(
					"expected package mapping for %s",
					tt.manager,
				)
			}

			if len(actual) != len(tt.expected) {
				t.Fatalf(
					"expected %d packages, got %d",
					len(tt.expected),
					len(actual),
				)
			}

			for i, expected := range tt.expected {
				if actual[i] != expected {
					t.Fatalf(
						"package %d: expected %q, got %q",
						i,
						expected,
						actual[i],
					)
				}
			}
		})
	}
}

func TestListProfiles(t *testing.T) {
	profiles := List()

	if len(profiles) != 3 {
		t.Fatalf(
			"expected 3 profiles, got %d",
			len(profiles),
		)
	}

	expected := []string{
		"essentials",
		"development",
		"multimedia",
	}

	for i, profile := range profiles {
		if profile.Name != expected[i] {
			t.Fatalf(
				"profile %d: expected %q, got %q",
				i,
				expected[i],
				profile.Name,
			)
		}
	}
}
func TestPackagesForUnknownManager(t *testing.T) {
	actual, ok := PackagesFor(Development, packages.Unknown)

	if ok {
		t.Fatal("expected no package mapping for unknown manager")
	}

	if len(actual) != 0 {
		t.Fatalf(
			"expected no packages for unknown manager, got %v",
			actual,
		)
	}
}
