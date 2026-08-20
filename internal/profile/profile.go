package profile

type Profile struct {
	Name  		 string
	Description  string
	Packages   []string
}

var Essentials = Profile{
	Name: "essentials",
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