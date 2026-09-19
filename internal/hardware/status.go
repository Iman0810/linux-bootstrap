package hardware

type HardwareStatus struct {
	GPUs        []GPU
	Nvidia      *NvidiaStatus
	NvidiaFound bool
}

func DetectHardware() HardwareStatus {
	return DetectHardwareWithRunner(OSCommandRunner{})
}

func DetectHardwareWithRunner(runner CommandRunner) HardwareStatus {
	gpus := DetectGPUsWithRunner(runner)

	status := HardwareStatus{
		GPUs: gpus,
	}

	for _, gpu := range gpus {
		if gpu.Vendor == NVIDIA {
			nvidia := DetectNvidiaDriverWithRunner(gpus, runner)

			status.Nvidia = &nvidia
			status.NvidiaFound = true

			break
		}
	}

	return status
}
