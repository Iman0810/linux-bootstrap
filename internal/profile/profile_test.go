package profile

import "testing"

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
		expected []string
	}{
		{
			name:    "Essentials",
			profile: Essentials,
			expected: []string{
				"git",
				"curl",
				"wget",
				"unzip",
			},
		},
		{
			name:    "Development",
			profile: Development,
			expected: []string{
				"git",
				"curl",
				"wget",
				"unzip",
				"build-essential",
			},
		},
		{
			name:    "Multimedia",
			profile: Multimedia,
			expected: []string{
				"ffmpeg",
				"vlc",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.profile.Packages) != len(tt.expected) {
				t.Fatalf(
					"expected %d packages, got %d",
					len(tt.expected),
					len(tt.profile.Packages),
				)
			}

			for i, expected := range tt.expected {
				if tt.profile.Packages[i] != expected {
					t.Fatalf(
						"package %d: expected %q, got %q",
						i,
						expected,
						tt.profile.Packages[i],
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
