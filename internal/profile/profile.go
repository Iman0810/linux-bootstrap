package profile

import "github.com/Iman0810/linux-bootstrap/internal/packages"

type Profile struct {
	Name        string
	Description string
	Packages    map[packages.Manager][]string
}

var Essentials = Profile{
	Name:        "essentials",
	Description: "Basic packages useful on a fresh Linux installation",
	Packages: map[packages.Manager][]string{
		packages.APT: {
			"git",
			"curl",
			"wget",
			"unzip",
		},
		packages.DNF: {
			"git",
			"curl",
			"wget",
			"unzip",
		},
		packages.Pacman: {
			"git",
			"curl",
			"wget",
			"unzip",
		},
	},
}

var Development = Profile{
	Name:        "development",
	Description: "Common tools for software development",
	Packages: map[packages.Manager][]string{
		packages.APT: {
			"git",
			"curl",
			"wget",
			"unzip",
			"build-essential",
		},
		packages.DNF: {
			"git",
			"curl",
			"wget",
			"unzip",
			"gcc",
			"gcc-c++",
			"make",
		},
		packages.Pacman: {
			"git",
			"curl",
			"wget",
			"unzip",
			"base-devel",
		},
	},
}

var Multimedia = Profile{
	Name:        "multimedia",
	Description: "Common multimedia tools and codecs",
	Packages: map[packages.Manager][]string{
		packages.APT: {
			"ffmpeg",
			"vlc",
		},
		packages.DNF: {
			"ffmpeg",
			"vlc",
		},
		packages.Pacman: {
			"ffmpeg",
			"vlc",
		},
	},
}

var profiles = []Profile{
	Essentials,
	Development,
	Multimedia,
}

func Get(name string) (Profile, bool) {
	for _, p := range profiles {
		if p.Name == name {
			return p, true
		}
	}

	return Profile{}, false
}

func List() []Profile {
	return profiles
}
func PackagesFor(p Profile, manager packages.Manager) ([]string, bool) {
	packages, ok := p.Packages[manager]

	return packages, ok
}
