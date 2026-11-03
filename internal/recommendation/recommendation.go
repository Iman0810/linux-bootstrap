package recommendation

import (
	"fmt"

	"github.com/Iman0810/linux-bootstrap/internal/doctor"
)

type Recommendation struct {
	Title       string
	Description string
	Command     string
}

func Generate(report doctor.Report) []Recommendation {
	var recommendations []Recommendation
	// Package manager availability
	if !report.PackageManagerAvailable {
		description := report.PackageManagerError

		if description == "" {
			description = fmt.Sprintf(
				"The expected package manager %s is unavailable.",
				report.PackageManager,
			)
		}

		recommendations = append(recommendations, Recommendation{
			Title:       "Package manager unavailable",
			Description: description,
			Command:     "Install or restore the expected package manager, then rerun linux-bootstrap doctor.",
		})
	}

	// NVIDIA driver
	if report.NvidiaFound && !report.NvidiaInstalled {
		recommendations = append(recommendations, Recommendation{
			Title:       "NVIDIA driver is missing",
			Description: "An NVIDIA GPU was detected, but no NVIDIA driver was detected.",
			Command:     "Install the appropriate NVIDIA driver for this system.",
		})
	}

	// Profiles
	for _, p := range report.Profiles {
		if p.Ready {
			continue
		}

		recommendations = append(recommendations, Recommendation{
			Title: fmt.Sprintf("Install %s profile", p.Name),
			Description: fmt.Sprintf(
				"Missing packages: %v",
				p.Missing,
			),
			Command: fmt.Sprintf(
				"linux-bootstrap setup --profile %s",
				p.Name,
			),
		})
	}

	return recommendations
}
