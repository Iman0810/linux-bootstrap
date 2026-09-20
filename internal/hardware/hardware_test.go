package hardware

import (
	"fmt"
	"testing"
)

type fakeCommandRunner struct {
	outputs map[string]string
	errs    map[string]error
}

func (f fakeCommandRunner) Output(command string, args ...string) (string, error) {
	if err, ok := f.errs[command]; ok {
		return "", err
	}

	return f.outputs[command], nil
}

func TestParseGPUs(t *testing.T) {
	input := `
0000:00:02.0 VGA compatible controller: Intel Corporation Raptor Lake-P [UHD Graphics]
0000:01:00.0 VGA compatible controller: NVIDIA Corporation GN20-P0-R-K2 [GeForce RTX 3050 6GB Laptop GPU]
0000:02:00.0 3D controller: Advanced Micro Devices, Inc. [AMD Radeon]
0000:03:00.0 Ethernet controller: Intel Corporation Ethernet Controller
`

	got := ParseGPUs(input)

	if len(got) != 3 {
		t.Fatalf("expected 3 GPUs, got %d", len(got))
	}

	tests := []struct {
		vendor GPUVendor
	}{
		{Intel},
		{NVIDIA},
		{AMD},
	}

	for i, tt := range tests {
		if got[i].Vendor != tt.vendor {
			t.Errorf(
				"GPU %d: expected vendor %q, got %q",
				i,
				tt.vendor,
				got[i].Vendor,
			)
		}
	}
}

func TestParseGPUsIgnoresNonGPUDevices(t *testing.T) {
	input := `
0000:00:14.0 USB controller: Intel Corporation
0000:00:1f.3 Audio device: Intel Corporation
0000:00:1f.6 Ethernet controller: Intel Corporation
`

	got := ParseGPUs(input)

	if len(got) != 0 {
		t.Fatalf("expected 0 GPUs, got %d", len(got))
	}
}

func TestParseGPUsEmptyOutput(t *testing.T) {
	got := ParseGPUs("")

	if len(got) != 0 {
		t.Fatalf("expected 0 GPUs, got %d", len(got))
	}
}

func TestParseNvidiaDriver(t *testing.T) {
	got := ParseNvidiaDriver("595.84\n")

	if !got.Installed {
		t.Fatal("expected NVIDIA driver to be installed")
	}

	if got.Version != "595.84" {
		t.Fatalf("expected version 595.84, got %q", got.Version)
	}
}

func TestParseNvidiaDriverEmptyOutput(t *testing.T) {
	got := ParseNvidiaDriver("")

	if got.Installed {
		t.Fatal("expected NVIDIA driver to be not installed")
	}

	if got.Version != "" {
		t.Fatalf("expected empty version, got %q", got.Version)
	}
}

func TestParseNvidiaDriverWhitespace(t *testing.T) {
	got := ParseNvidiaDriver("   \n\t")

	if got.Installed {
		t.Fatal("expected NVIDIA driver to be not installed")
	}
}

func TestParseNvidiaDriverTrimsVersion(t *testing.T) {
	got := ParseNvidiaDriver("  595.84  \n")

	if !got.Installed {
		t.Fatal("expected NVIDIA driver to be installed")
	}

	if got.Version != "595.84" {
		t.Fatalf("expected trimmed version 595.84, got %q", got.Version)
	}
}
func TestDetectGPUsWithRunner(t *testing.T) {
	runner := fakeCommandRunner{
		outputs: map[string]string{
			"lspci": `
	0000:00:02.0 VGA compatible controller: Intel Corporation UHD Graphics
	0000:01:00.0 VGA compatible controller: NVIDIA Corporation GeForce RTX 3050
	`,
		},
	}

	got := DetectGPUsWithRunner(runner)

	if len(got) != 2 {
		t.Fatalf("expected 2 GPUs, got %d", len(got))
	}

	if got[0].Vendor != Intel {
		t.Fatalf("expected first GPU to be Intel, got %q", got[0].Vendor)
	}

	if got[1].Vendor != NVIDIA {
		t.Fatalf("expected second GPU to be NVIDIA, got %q", got[1].Vendor)
	}
}

func TestDetectGPUsWithRunnerError(t *testing.T) {
	runner := fakeCommandRunner{
		errs: map[string]error{
			"lspci": fmt.Errorf("lspci failed"),
		},
	}
	got := DetectGPUsWithRunner(runner)

	if len(got) != 0 {
		t.Fatalf("expected 0 GPUs on error, got %d", len(got))
	}
}
func TestDetectNvidiaDriverWithRunner(t *testing.T) {
	gpus := []GPU{
		{
			Vendor: NVIDIA,
			Name:   "NVIDIA GeForce RTX 3050",
		},
	}

	runner := fakeCommandRunner{
		outputs: map[string]string{
			"nvidia-smi": "595.84\n",
		},
	}
	got := DetectNvidiaDriverWithRunner(gpus, runner)

	if !got.Installed {
		t.Fatal("expected NVIDIA driver to be installed")
	}

	if got.Version != "595.84" {
		t.Fatalf("expected version 595.84, got %q", got.Version)
	}
}

func TestDetectNvidiaDriverWithRunnerNoNvidia(t *testing.T) {
	gpus := []GPU{
		{
			Vendor: Intel,
			Name:   "Intel UHD Graphics",
		},
	}

	runner := fakeCommandRunner{
		outputs: map[string]string{
			"nvidia-smi": "595.84\n",
		},
	}
	got := DetectNvidiaDriverWithRunner(gpus, runner)

	if got.Installed {
		t.Fatal("expected NVIDIA driver to be not installed")
	}

	if got.Version != "" {
		t.Fatalf("expected empty version, got %q", got.Version)
	}
}

func TestDetectNvidiaDriverWithRunnerError(t *testing.T) {
	gpus := []GPU{
		{
			Vendor: NVIDIA,
			Name:   "NVIDIA GeForce RTX 3050",
		},
	}

	runner := fakeCommandRunner{
		errs: map[string]error{
			"nvidia-smi": fmt.Errorf("nvidia-smi failed"),
		},
	}

	got := DetectNvidiaDriverWithRunner(gpus, runner)

	if got.Installed {
		t.Fatal("expected NVIDIA driver to be not installed")
	}
}
func TestDetectHardwareWithRunner(t *testing.T) {
	runner := fakeCommandRunner{
		outputs: map[string]string{
			"lspci": `
	0000:00:02.0 VGA compatible controller: Intel Corporation UHD Graphics
	0000:01:00.0 VGA compatible controller: NVIDIA Corporation GeForce RTX 3050
	`,
			"nvidia-smi": "595.84\n",
		},
	}

	got := DetectHardwareWithRunner(runner)

	if len(got.GPUs) != 2 {
		t.Fatalf("expected 2 GPUs, got %d", len(got.GPUs))
	}

	if !got.NvidiaFound {
		t.Fatal("expected NVIDIA to be found")
	}

	if got.Nvidia == nil {
		t.Fatal("expected NVIDIA status")
	}

	if !got.Nvidia.Installed {
		t.Fatal("expected NVIDIA driver to be installed")
	}

	if got.Nvidia.Version != "595.84" {
		t.Fatalf("expected NVIDIA version 595.84, got %q", got.Nvidia.Version)
	}
}

func TestDetectHardwareWithRunnerNoNvidia(t *testing.T) {
	runner := fakeCommandRunner{
		outputs: map[string]string{
			"lspci": `
	0000:00:02.0 VGA compatible controller: Intel Corporation UHD Graphics
	`,
		},
	}

	got := DetectHardwareWithRunner(runner)

	if len(got.GPUs) != 1 {
		t.Fatalf("expected 1 GPU, got %d", len(got.GPUs))
	}

	if got.NvidiaFound {
		t.Fatal("expected NVIDIA not to be found")
	}

	if got.Nvidia != nil {
		t.Fatal("expected NVIDIA status to be nil")
	}
}

func TestDetectHardwareWithRunnerNvidiaDriverFailure(t *testing.T) {
	
}
