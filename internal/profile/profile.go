package profile

type Profile struct {
	Name        string
	Description string
	Packages    []string
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
