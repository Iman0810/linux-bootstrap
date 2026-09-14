package system

import (
	"bufio"
	"os"
	"strings"
)

type OSInfo struct {
	Name    string
	Version string
	ID      string
	IDLike  string
}

func GetOSInfo() (OSInfo, error) {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return OSInfo{}, err
	}
	defer file.Close()

	var builder strings.Builder

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		builder.WriteString(scanner.Text())
		builder.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return OSInfo{}, err
	}

	return ParseOSRelease(builder.String()), nil
}

func ParseOSRelease(content string) OSInfo {
	info := OSInfo{}

	lines := strings.Split(content, "\n")

	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := strings.Trim(parts[1], `"`)

		switch key {
		case "NAME":
			info.Name = value
		case "VERSION":
			info.Version = value
		case "ID":
			info.ID = value
		case "ID_LIKE":
			info.IDLike = value
		}
	}

	return info
}
