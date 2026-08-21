package hardware

import (
	"os/exec"
	"strings"
)

type NvidiaStatus struct {
	Installed bool
	Version   string
}

func DetectNvidiaDriver() NvidiaStatus {
	cmd := exec.Command("nvidia-smi", "--query-gpu=driver_version", "--format=csv,noheader")

	output, err := cmd.Output()
	if err != nil {
		return NvidiaStatus{
			Installed: false,
		}
	}

	version := strings.TrimSpace(string(output))

	if version == "" {
		return NvidiaStatus{
			Installed: false,
		}
	}

	return NvidiaStatus{
		Installed: true,
		Version:   version,
	}
}
