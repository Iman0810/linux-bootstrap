package hardware

import "testing"

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
