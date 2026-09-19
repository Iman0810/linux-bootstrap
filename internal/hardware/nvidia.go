package hardware

import "strings"

type NvidiaStatus struct {
	Installed bool
	Version   string
}

func DetectNvidiaDriver(gpus []GPU) NvidiaStatus {
	return DetectNvidiaDriverWithRunner(gpus, OSCommandRunner{})
}

func DetectNvidiaDriverWithRunner(
	gpus []GPU,
	runner CommandRunner,
) NvidiaStatus {
	hasNvidia := false

	for _, gpu := range gpus {
		if gpu.Vendor == NVIDIA {
			hasNvidia = true
			break
		}
	}

	if !hasNvidia {
		return NvidiaStatus{}
	}

	output, err := runner.Output(
		"nvidia-smi",
		"--query-gpu=driver_version",
		"--format=csv,noheader",
	)

	if err != nil {
		return NvidiaStatus{
			Installed: false,
		}
	}

	return ParseNvidiaDriver(output)
}

func ParseNvidiaDriver(output string) NvidiaStatus {
	version := strings.TrimSpace(output)

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
