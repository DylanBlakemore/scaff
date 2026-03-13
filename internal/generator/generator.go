package generator

type Request struct {
	Name       string
	ModulePath string
	Style      string
	Features   []string
	OutputDir  string
	License    string
}

type Result struct {
	Files []string
}

type ProjectGenerator interface {
	Generate(req Request) (*Result, error)
}
