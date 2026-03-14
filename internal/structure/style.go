package structure

type Style struct {
	Name  string
	Paths map[string]string
}

func MinimalPackageStyle() Style {
	return Style{
		Name: "minimal",
		Paths: map[string]string{
			"root":   ".",
			"source": ".",
			"test":   ".",
			"ci":     ".github/workflows",
			"meta":   ".",
		},
	}
}

func (s Style) Resolve(component string) string {
	if p, ok := s.Paths[component]; ok {
		return p
	}
	return "."
}
