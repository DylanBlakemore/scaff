package structure

type Style struct {
	Name  string
	Paths map[string]string
}

func DefaultPackageStyle() Style {
	return Style{
		Name: "default",
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
