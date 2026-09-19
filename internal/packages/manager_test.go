package packages

import (
	"fmt"
	"testing"

	"github.com/Iman0810/linux-bootstrap/internal/system"
)

type fakeRunner struct {
	lastCommand string
	lastArgs    []string

	runErr    error
	output    string
	outputErr error
}

func (f *fakeRunner) Run(command string, args ...string) error {
	f.lastCommand = command
	f.lastArgs = args
	return f.runErr
}

func (f *fakeRunner) Output(command string, args ...string) (string, error) {
	f.lastCommand = command
	f.lastArgs = args
	return f.output, f.outputErr
}

func TestDetectManager(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected Manager
	}{
		{
			name:     "Ubuntu",
			id:       "ubuntu",
			expected: APT,
		},
		{
			name:     "Debian",
			id:       "debian",
			expected: APT,
		},
		{
			name:     "Pop OS",
			id:       "pop",
			expected: APT,
		},
		{
			name:     "Fedora",
			id:       "fedora",
			expected: DNF,
		},
		{
			name:     "RHEL",
			id:       "rhel",
			expected: DNF,
		},
		{
			name:     "Arch",
			id:       "arch",
			expected: Pacman,
		},
		{
			name:     "Manjaro",
			id:       "manjaro",
			expected: Pacman,
		},
		{
			name:     "Unknown",
			id:       "something-random",
			expected: Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			osInfo := system.OSInfo{
				ID: tt.id,
			}

			got := DetectManager(osInfo)

			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
func TestAptManagerUpdate(t *testing.T) {
	fake := &fakeRunner{}
	manager := AptManager{Runner: fake}

	err := manager.Update()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if fake.lastCommand != "sudo" {
		t.Fatalf("expected command sudo, got %q", fake.lastCommand)
	}

	expectedArgs := []string{"apt", "update"}

	if len(fake.lastArgs) != len(expectedArgs) {
		t.Fatalf(
			"expected %d args, got %d",
			len(expectedArgs),
			len(fake.lastArgs),
		)
	}

	for i, expected := range expectedArgs {
		if fake.lastArgs[i] != expected {
			t.Fatalf(
				"expected arg %d to be %q, got %q",
				i,
				expected,
				fake.lastArgs[i],
			)
		}
	}
}
func TestAptManagerUpdateError(t *testing.T) {
	expectedErr := fmt.Errorf("apt update failed")

	fake := &fakeRunner{
		runErr: expectedErr,
	}

	manager := AptManager{Runner: fake}

	err := manager.Update()
	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}
func TestAptManagerInstall(t *testing.T) {
	fake := &fakeRunner{}
	manager := AptManager{Runner: fake}

	err := manager.Install("git", "curl", "wget")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if fake.lastCommand != "sudo" {
		t.Fatalf("expected command sudo, got %q", fake.lastCommand)
	}

	expectedArgs := []string{
		"apt",
		"install",
		"-y",
		"git",
		"curl",
		"wget",
	}

	if len(fake.lastArgs) != len(expectedArgs) {
		t.Fatalf(
			"expected %d args, got %d",
			len(expectedArgs),
			len(fake.lastArgs),
		)
	}

	for i, expected := range expectedArgs {
		if fake.lastArgs[i] != expected {
			t.Fatalf(
				"expected arg %d to be %q, got %q",
				i,
				expected,
				fake.lastArgs[i],
			)
		}
	}
}
func TestAptManagerIsInstalled(t *testing.T) {
	fake := &fakeRunner{
		output: "Package: git\nStatus: install ok installed\n",
	}

	manager := AptManager{Runner: fake}

	if !manager.IsInstalled("git") {
		t.Fatal("expected package to be installed")
	}

	if fake.lastCommand != "dpkg" {
		t.Fatalf("expected command dpkg, got %q", fake.lastCommand)
	}

	expectedArgs := []string{"-s", "git"}

	if len(fake.lastArgs) != len(expectedArgs) {
		t.Fatalf(
			"expected %d args, got %d",
			len(expectedArgs),
			len(fake.lastArgs),
		)
	}

	for i, expected := range expectedArgs {
		if fake.lastArgs[i] != expected {
			t.Fatalf(
				"expected arg %d to be %q, got %q",
				i,
				expected,
				fake.lastArgs[i],
			)
		}
	}
}

func TestAptManagerIsNotInstalled(t *testing.T) {
	fake := &fakeRunner{
		outputErr: fmt.Errorf("package not installed"),
	}

	manager := AptManager{Runner: fake}

	if manager.IsInstalled("some-package") {
		t.Fatal("expected package to not be installed")
	}
}
func TestDnfManagerUpdate(t *testing.T) {
	fake := &fakeRunner{}
	manager := DnfManager{Runner: fake}

	err := manager.Update()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if fake.lastCommand != "sudo" {
		t.Fatalf("expected command sudo, got %q", fake.lastCommand)
	}

	expectedArgs := []string{"dnf", "makecache"}

	if len(fake.lastArgs) != len(expectedArgs) {
		t.Fatalf(
			"expected %d args, got %d",
			len(expectedArgs),
			len(fake.lastArgs),
		)
	}

	for i, expected := range expectedArgs {
		if fake.lastArgs[i] != expected {
			t.Fatalf(
				"expected arg %d to be %q, got %q",
				i,
				expected,
				fake.lastArgs[i],
			)
		}
	}
}

func TestDnfManagerInstall(t *testing.T) {
	fake := &fakeRunner{}
	manager := DnfManager{Runner: fake}

	err := manager.Install("git", "curl")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if fake.lastCommand != "sudo" {
		t.Fatalf("expected command sudo, got %q", fake.lastCommand)
	}

	expectedArgs := []string{
		"dnf",
		"install",
		"-y",
		"git",
		"curl",
	}

	if len(fake.lastArgs) != len(expectedArgs) {
		t.Fatalf(
			"expected %d args, got %d",
			len(expectedArgs),
			len(fake.lastArgs),
		)
	}

	for i, expected := range expectedArgs {
		if fake.lastArgs[i] != expected {
			t.Fatalf(
				"expected arg %d to be %q, got %q",
				i,
				expected,
				fake.lastArgs[i],
			)
		}
	}
}

func TestDnfManagerIsInstalled(t *testing.T) {
	fake := &fakeRunner{
		output: "git-2.47.1-1.fc44.x86_64",
	}

	manager := DnfManager{Runner: fake}

	if !manager.IsInstalled("git") {
		t.Fatal("expected package to be installed")
	}

	if fake.lastCommand != "rpm" {
		t.Fatalf("expected command rpm, got %q", fake.lastCommand)
	}

	expectedArgs := []string{"-q", "git"}

	if len(fake.lastArgs) != len(expectedArgs) {
		t.Fatalf(
			"expected %d args, got %d",
			len(expectedArgs),
			len(fake.lastArgs),
		)
	}

	for i, expected := range expectedArgs {
		if fake.lastArgs[i] != expected {
			t.Fatalf(
				"expected arg %d to be %q, got %q",
				i,
				expected,
				fake.lastArgs[i],
			)
		}
	}
}
func TestDnfManagerIsNotInstalled(t *testing.T) {
	fake := &fakeRunner{
		outputErr: fmt.Errorf("package not installed"),
	}

	manager := DnfManager{Runner: fake}

	if manager.IsInstalled("some-package") {
		t.Fatal("expected package to not be installed")
	}
}
func TestPacmanManagerUpdate(t *testing.T) {
	fake := &fakeRunner{}
	manager := PacmanManager{Runner: fake}

	err := manager.Update()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if fake.lastCommand != "sudo" {
		t.Fatalf("expected command sudo, got %q", fake.lastCommand)
	}

	expectedArgs := []string{"pacman", "-Syu"}

	if len(fake.lastArgs) != len(expectedArgs) {
		t.Fatalf(
			"expected %d args, got %d",
			len(expectedArgs),
			len(fake.lastArgs),
		)
	}

	for i, expected := range expectedArgs {
		if fake.lastArgs[i] != expected {
			t.Fatalf(
				"expected arg %d to be %q, got %q",
				i,
				expected,
				fake.lastArgs[i],
			)
		}
	}
}

func TestPacmanManagerInstall(t *testing.T) {
	fake := &fakeRunner{}
	manager := PacmanManager{Runner: fake}

	err := manager.Install("git", "curl")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if fake.lastCommand != "sudo" {
		t.Fatalf("expected command sudo, got %q", fake.lastCommand)
	}

	expectedArgs := []string{
		"pacman",
		"-S",
		"--noconfirm",
		"git",
		"curl",
	}

	if len(fake.lastArgs) != len(expectedArgs) {
		t.Fatalf(
			"expected %d args, got %d",
			len(expectedArgs),
			len(fake.lastArgs),
		)
	}

	for i, expected := range expectedArgs {
		if fake.lastArgs[i] != expected {
			t.Fatalf(
				"expected arg %d to be %q, got %q",
				i,
				expected,
				fake.lastArgs[i],
			)
		}
	}
}

func TestPacmanManagerIsInstalled(t *testing.T) {
	fake := &fakeRunner{
		output: "git 2.50.1-1",
	}

	manager := PacmanManager{Runner: fake}

	if !manager.IsInstalled("git") {
		t.Fatal("expected package to be installed")
	}

	if fake.lastCommand != "pacman" {
		t.Fatalf("expected command pacman, got %q", fake.lastCommand)
	}

	expectedArgs := []string{"-Q", "git"}

	if len(fake.lastArgs) != len(expectedArgs) {
		t.Fatalf(
			"expected %d args, got %d",
			len(expectedArgs),
			len(fake.lastArgs),
		)
	}

	for i, expected := range expectedArgs {
		if fake.lastArgs[i] != expected {
			t.Fatalf(
				"expected arg %d to be %q, got %q",
				i,
				expected,
				fake.lastArgs[i],
			)
		}
	}
}

func TestPacmanManagerIsNotInstalled(t *testing.T) {
	fake := &fakeRunner{
		outputErr: fmt.Errorf("package not installed"),
	}

	manager := PacmanManager{Runner: fake}

	if manager.IsInstalled("some-package") {
		t.Fatal("expected package to not be installed")
	}
}
