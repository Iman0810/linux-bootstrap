package recommendation

import (
	"reflect"
	"testing"

	"github.com/Iman0810/linux-bootstrap/internal/doctor"
)

func TestGenerateNoRecommendations(t *testing.T) {
	report := doctor.Report{}

	got := Generate(report)

	if len(got) != 0 {
		t.Fatalf("expected 0 recommendations, got %d", len(got))
	}
}

func TestGenerateNvidiaRecommendation(t *testing.T) {
	report := doctor.Report{
		NvidiaFound:     true,
		NvidiaInstalled: false,
	}

	got := Generate(report)

	if len(got) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(got))
	}

	expected := Recommendation{
		Title:       "NVIDIA driver is missing",
		Description: "An NVIDIA GPU was detected, but no NVIDIA driver was detected.",
		Command:     "Install the appropriate NVIDIA driver for this system.",
	}

	if !reflect.DeepEqual(got[0], expected) {
		t.Fatalf("unexpected recommendation:\n got:  %+v\nwant: %+v", got[0], expected)
	}
}

func TestGenerateSkipsInstalledNvidiaDriver(t *testing.T) {
	report := doctor.Report{
		NvidiaFound:     true,
		NvidiaInstalled: true,
	}

	got := Generate(report)

	if len(got) != 0 {
		t.Fatalf("expected 0 recommendations, got %d", len(got))
	}
}

func TestGenerateProfileRecommendation(t *testing.T) {
	report := doctor.Report{
		Profiles: []doctor.ProfileStatus{
			{
				Name:    "multimedia",
				Ready:   false,
				Missing: []string{"ffmpeg", "vlc"},
			},
		},
	}

	got := Generate(report)

	if len(got) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(got))
	}

	expected := Recommendation{
		Title:       "Install multimedia profile",
		Description: "Missing packages: [ffmpeg vlc]",
		Command:     "linux-bootstrap setup --profile multimedia",
	}

	if !reflect.DeepEqual(got[0], expected) {
		t.Fatalf("unexpected recommendation:\n got:  %+v\nwant: %+v", got[0], expected)
	}
}

func TestGenerateSkipsReadyProfiles(t *testing.T) {
	report := doctor.Report{
		Profiles: []doctor.ProfileStatus{
			{
				Name:    "essentials",
				Ready:   true,
				Missing: nil,
			},
		},
	}

	got := Generate(report)

	if len(got) != 0 {
		t.Fatalf("expected 0 recommendations, got %d", len(got))
	}
}

func TestGenerateMultipleRecommendations(t *testing.T) {
	report := doctor.Report{
		NvidiaFound:     true,
		NvidiaInstalled: false,
		Profiles: []doctor.ProfileStatus{
			{
				Name:    "multimedia",
				Ready:   false,
				Missing: []string{"ffmpeg", "vlc"},
			},
			{
				Name:  "development",
				Ready: true,
			},
			{
				Name:    "essentials",
				Ready:   false,
				Missing: []string{"git"},
			},
		},
	}

	got := Generate(report)

	if len(got) != 3 {
		t.Fatalf("expected 3 recommendations, got %d", len(got))
	}

	if got[0].Title != "NVIDIA driver is missing" {
		t.Errorf("unexpected first recommendation: %q", got[0].Title)
	}

	if got[1].Title != "Install multimedia profile" {
		t.Errorf("unexpected second recommendation: %q", got[1].Title)
	}

	if got[2].Title != "Install essentials profile" {
		t.Errorf("unexpected third recommendation: %q", got[2].Title)
	}
}
