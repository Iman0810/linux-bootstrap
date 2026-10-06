package system

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestParseOSReleasePopOS(t *testing.T) {
	input := `NAME="Pop!_OS"
VERSION="24.04 LTS"
ID=pop
ID_LIKE="ubuntu debian"
`

	got := ParseOSRelease(input)

	if got.Name != "Pop!_OS" {
		t.Errorf("expected name %q, got %q", "Pop!_OS", got.Name)
	}

	if got.Version != "24.04 LTS" {
		t.Errorf("expected version %q, got %q", "24.04 LTS", got.Version)
	}

	if got.ID != "pop" {
		t.Errorf("expected ID %q, got %q", "pop", got.ID)
	}

	if got.IDLike != "ubuntu debian" {
		t.Errorf("expected ID_LIKE %q, got %q", "ubuntu debian", got.IDLike)
	}
}

func TestParseOSReleaseUbuntu(t *testing.T) {
	input := `NAME="Ubuntu"
VERSION="24.04.3 LTS (Noble Numbat)"
ID=ubuntu
ID_LIKE=debian
`

	got := ParseOSRelease(input)

	if got.Name != "Ubuntu" {
		t.Errorf("expected name %q, got %q", "Ubuntu", got.Name)
	}

	if got.Version != "24.04.3 LTS (Noble Numbat)" {
		t.Errorf("unexpected version: %q", got.Version)
	}

	if got.ID != "ubuntu" {
		t.Errorf("expected ID %q, got %q", "ubuntu", got.ID)
	}

	if got.IDLike != "debian" {
		t.Errorf("expected ID_LIKE %q, got %q", "debian", got.IDLike)
	}
}

func TestParseOSReleaseFedora(t *testing.T) {
	input := `NAME="Fedora Linux"
VERSION="44 (Workstation Edition)"
ID=fedora
VERSION_ID=44
`

	got := ParseOSRelease(input)

	if got.Name != "Fedora Linux" {
		t.Errorf("expected name %q, got %q", "Fedora Linux", got.Name)
	}

	if got.Version != "44 (Workstation Edition)" {
		t.Errorf("unexpected version: %q", got.Version)
	}

	if got.ID != "fedora" {
		t.Errorf("expected ID %q, got %q", "fedora", got.ID)
	}
}

func TestParseOSReleaseArch(t *testing.T) {
	input := `NAME="Arch Linux"
PRETTY_NAME="Arch Linux"
ID=arch
BUILD_ID=rolling
`

	got := ParseOSRelease(input)

	if got.Name != "Arch Linux" {
		t.Errorf("expected name %q, got %q", "Arch Linux", got.Name)
	}

	if got.ID != "arch" {
		t.Errorf("expected ID %q, got %q", "arch", got.ID)
	}
}

func TestParseOSReleaseIgnoresUnknownFields(t *testing.T) {
	input := `NAME="Test Linux"
VERSION="1.0"
ID=test
UNKNOWN_FIELD="ignored"
PRETTY_NAME="Test Linux 1.0"
`

	got := ParseOSRelease(input)

	expected := OSInfo{
		Name:    "Test Linux",
		Version: "1.0",
		ID:      "test",
	}

	if got != expected {
		t.Fatalf("unexpected OS info: got %+v, want %+v", got, expected)
	}
}

func TestParseOSReleaseIgnoresMalformedLines(t *testing.T) {
	input := `NAME="Test Linux"
this line has no equals sign
VERSION="1.0"
ID=test
another invalid line
`

	got := ParseOSRelease(input)

	expected := OSInfo{
		Name:    "Test Linux",
		Version: "1.0",
		ID:      "test",
	}

	if got != expected {
		t.Fatalf("unexpected OS info: got %+v, want %+v", got, expected)
	}
}

func TestParseOSReleaseEmptyInput(t *testing.T) {
	got := ParseOSRelease("")

	if got != (OSInfo{}) {
		t.Fatalf("expected empty OSInfo, got %+v", got)
	}
}
func TestParseOSReleaseReader(t *testing.T) {
	input := `NAME="Test Linux"
VERSION="1.0"
ID=test
ID_LIKE=debian
`

	got, err := parseOSReleaseReader(strings.NewReader(input))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := OSInfo{
		Name:    "Test Linux",
		Version: "1.0",
		ID:      "test",
		IDLike:  "debian",
	}

	if got != expected {
		t.Fatalf("unexpected OS info: got %+v, want %+v", got, expected)
	}
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) {
	return 0, fmt.Errorf("read failed")
}

func TestParseOSReleaseReaderError(t *testing.T) {
	_, err := parseOSReleaseReader(errorReader{})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "read failed" {
		t.Fatalf("expected read failed error, got %v", err)
	}
}
func TestGetOSInfo(t *testing.T) {
	originalOpen := openOSRelease
	defer func() {
		openOSRelease = originalOpen
	}()

	tempFile, err := os.CreateTemp("", "os-release-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = tempFile.WriteString(`NAME="Test Linux"
VERSION="1.0"
ID=test
ID_LIKE=debian
`)
	if err != nil {
		t.Fatalf("failed to write test data: %v", err)
	}

	if _, err := tempFile.Seek(0, 0); err != nil {
		t.Fatalf("failed to seek test file: %v", err)
	}

	openOSRelease = func(_ string) (*os.File, error) {
		return os.Open(tempFile.Name())
	}

	got, err := GetOSInfo()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := OSInfo{
		Name:    "Test Linux",
		Version: "1.0",
		ID:      "test",
		IDLike:  "debian",
	}

	if got != expected {
		t.Fatalf("unexpected OS info: got %+v, want %+v", got, expected)
	}
}
func TestGetOSInfoOpenFailure(t *testing.T) {
	originalOpen := openOSRelease
	defer func() {
		openOSRelease = originalOpen
	}()

	expectedErr := fmt.Errorf("test open failure")

	openOSRelease = func(_ string) (*os.File, error) {
		return nil, expectedErr
	}

	got, err := GetOSInfo()

	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if got != (OSInfo{}) {
		t.Fatalf("expected empty OSInfo, got %+v", got)
	}
}
