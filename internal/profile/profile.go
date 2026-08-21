package profile

type Profile struct {
	Name        string
	Description string
	Packages    []string
}

func Get(name string) (Profile, bool) {
	switch name {
	case "essentials":
		return Essentials, true

	case "development":
		return Development, true

	case "multimedia":
		return Multimedia, true

	default:
		return Profile{}, false
	}
}

var Essentials = Profile{
	Name:        "essentials",
	Description: "Basic packages useful on a fresh Linux installation",
	Packages: []string{
		"git",
		"curl",
		"wget",
		"unzip",
	},
}

var Development = Profile{
	Name:        "development",
	Description: "Common tools for software development",
	Packages: []string{
		"git",
		"curl",
		"wget",
		"unzip",
		"build-essential",
	},
}

var Multimedia = Profile{
	Name:        "multimedia",
	Description: "Common multimedia tools and codecs",
	Packages: []string{
		"ffmpeg",
		"vlc",
	},
}