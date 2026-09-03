 package hardware

import (
	"os/exec"
	"strings"
)

type GPUVendor string

const (
	NVIDIA  GPUVendor = "nvidia"
	AMD     GPUVendor = "amd"
	Intel   GPUVendor = "intel"
	Unknown GPUVendor = "unknown"
)

type GPU struct {
	Vendor GPUVendor
	Name   string
}

func DetectGPUs() []GPU {
	cmd := exec.Command("lspci")

	output, err := cmd.Output()
	if err != nil {
		return []GPU{}
	}

	return ParseGPUs(string(output))
}

func ParseGPUs(output string) []GPU {
	lines := strings.Split(output, "\n")
	var gpus []GPU

	for _, line := range lines {
		lower := strings.ToLower(line)

		if !strings.Contains(lower, "vga") &&
			!strings.Contains(lower, "3d controller") {
			continue
		}

		switch {
		case strings.Contains(lower, "nvidia"):
			gpus = append(gpus, GPU{
				Vendor: NVIDIA,
				Name:   strings.TrimSpace(line),
			})

		case strings.Contains(lower, "amd"):
			gpus = append(gpus, GPU{
				Vendor: AMD,
				Name:   strings.TrimSpace(line),
			})

		case strings.Contains(lower, "intel"):
			gpus = append(gpus, GPU{
				Vendor: Intel,
				Name:   strings.TrimSpace(line),
			})
		}
	}

	return gpus
}

