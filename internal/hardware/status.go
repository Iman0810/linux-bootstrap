package hardware

type HardwareStatus struct {
	GPUs        []GPU
	Nvidia      *NvidiaStatus
	NvidiaFound bool
}

func DetectHardware() HardwareStatus {
	gpus := DetectGPUs()

	status := HardwareStatus{
		GPUs: gpus,
	}

	for _, gpu := range gpus {
		if gpu.Vendor == NVIDIA {
			nvidia := DetectNvidiaDriver(gpus)

			status.Nvidia = &nvidia
			status.NvidiaFound = true

			break
		}
	}

	return status
}
