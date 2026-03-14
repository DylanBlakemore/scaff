package pkg

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"scaff/internal/features"
	"scaff/internal/filesystem"
	"scaff/internal/generator"
	"scaff/internal/metadata"
	"scaff/internal/structure"
	"scaff/internal/templates"
)

type Deps struct {
	Style         structure.Style
	NewWriter     func(baseDir string) *filesystem.Writer
	Render        func(name string, data templates.Data) (string, error)
	ApplyFeatures func(w *filesystem.Writer, style structure.Style, data templates.Data, projectType string, featureNames []string) ([]string, error)
	GoVersion     func() string
	Year          func() int
}

type Generator struct {
	deps Deps
}

const licenseNone = "none"
const goVersion = "1.25.8"

func New() *Generator {
	return NewWithDeps(Deps{})
}

func NewWithDeps(d Deps) *Generator {
	if d.Style.Name == "" {
		d.Style = structure.MinimalPackageStyle()
	}
	if d.NewWriter == nil {
		d.NewWriter = filesystem.NewWriter
	}
	if d.Render == nil {
		d.Render = templates.Render
	}
	if d.ApplyFeatures == nil {
		d.ApplyFeatures = features.Apply
	}
	if d.GoVersion == nil {
		d.GoVersion = currentGoVersion
	}
	if d.Year == nil {
		d.Year = currentYear
	}

	return &Generator{deps: d}
}

func (g *Generator) Generate(req generator.Request) (*generator.Result, error) {
	if req.License == "" {
		req.License = licenseNone
	}
	if err := validate(req); err != nil {
		return nil, err
	}

	ctx, err := g.prepare(req)
	if err != nil {
		return nil, err
	}

	written, err := g.writeBaseFiles(ctx)
	if err != nil {
		return nil, err
	}

	licenseFiles, err := writeLicenseFiles(ctx.w, ctx.style, ctx.data, req.License)
	if err != nil {
		return nil, err
	}
	written = append(written, licenseFiles...)

	featureFiles, err := g.deps.ApplyFeatures(ctx.w, ctx.style, ctx.data, "package", req.Features)
	if err != nil {
		return nil, err
	}
	written = append(written, featureFiles...)

	metaPath, err := g.writeMetadata(ctx, req)
	if err != nil {
		return nil, err
	}
	written = append(written, metaPath)

	return &generator.Result{Files: written}, nil
}

type generateCtx struct {
	style structure.Style
	w     *filesystem.Writer
	data  templates.Data
}

func (g *Generator) prepare(req generator.Request) (*generateCtx, error) {
	style := g.deps.Style

	targetDir := req.OutputDir
	if targetDir == "" {
		targetDir = req.Name
	}

	w := g.deps.NewWriter(targetDir)
	if err := w.EnsureDir(); err != nil {
		return nil, err
	}

	data := templates.Data{
		Name:       sanitisePkgName(req.Name),
		ModulePath: req.ModulePath,
		GoVersion:  g.deps.GoVersion(),
		Year:       g.deps.Year(),
	}

	return &generateCtx{style: style, w: w, data: data}, nil
}

func (g *Generator) writeBaseFiles(ctx *generateCtx) ([]string, error) {
	files := []struct {
		component    string
		outputName   string
		templateName string
	}{
		{"root", ".gitignore", "gitignore"},
		{"root", ".editorconfig", "editorconfig"},
		{"root", ".golangci.yml", "golangci.yml"},
		{"root", "go.mod", "go.mod"},
		{"root", "CHANGELOG.md", "changelog.md"},
		{"root", "CONTRIBUTING.md", "contributing.md"},
		{"root", "SECURITY.md", "security.md"},
		{"source", ctx.data.Name + ".go", "package.go"},
		{"test", ctx.data.Name + "_test.go", "package_test.go"},
		{"root", "README.md", "readme.md"},
	}

	var written []string
	for _, f := range files {
		content, err := g.deps.Render(f.templateName, ctx.data)
		if err != nil {
			return nil, err
		}
		relPath := filepath.Join(ctx.style.Resolve(f.component), f.outputName)
		if err := ctx.w.WriteFile(relPath, content); err != nil {
			return nil, err
		}
		written = append(written, relPath)
	}
	return written, nil
}

func (g *Generator) writeMetadata(ctx *generateCtx, req generator.Request) (string, error) {
	meta := metadata.Project{
		ProjectType: "package",
		Style:       ctx.style.Name,
		Features:    req.Features,
	}
	metaContent, err := meta.Marshal()
	if err != nil {
		return "", err
	}
	metaPath := filepath.Join(ctx.style.Resolve("meta"), metadata.FileName)
	if err := ctx.w.WriteFile(metaPath, metaContent); err != nil {
		return "", err
	}
	return metaPath, nil
}

func validate(req generator.Request) error {
	if req.Name == "" {
		return fmt.Errorf("project name is required")
	}
	if req.ModulePath == "" {
		return fmt.Errorf("module path is required")
	}
	if req.Style != "minimal" {
		return fmt.Errorf("unsupported architecture style: %q (only \"minimal\" is available)", req.Style)
	}
	validFeatures := []string{"agents", "makefile", "ci"}
	for _, f := range req.Features {
		if !slices.Contains(validFeatures, f) {
			return fmt.Errorf("unknown feature: %q", f)
		}
	}

	if _, ok := licenseTemplate(req.License); !ok {
		return fmt.Errorf("unsupported license: %q", req.License)
	}
	return nil
}

func writeLicenseFiles(w *filesystem.Writer, style structure.Style, data templates.Data, license string) ([]string, error) {
	if license == licenseNone {
		content, err := templates.Render("copyright", data)
		if err != nil {
			return nil, err
		}
		relPath := filepath.Join(style.Resolve("root"), "COPYRIGHT")
		if err := w.WriteFile(relPath, content); err != nil {
			return nil, err
		}
		return []string{relPath}, nil
	}

	tmpl, _ := licenseTemplate(license)
	content, err := templates.Render(tmpl, data)
	if err != nil {
		return nil, err
	}
	relPath := filepath.Join(style.Resolve("root"), "LICENSE")
	if err := w.WriteFile(relPath, content); err != nil {
		return nil, err
	}
	return []string{relPath}, nil
}

func licenseTemplate(spdx string) (string, bool) {
	switch spdx {
	case licenseNone:
		return "", true
	case "MIT":
		return "licenses/mit", true
	case "Apache-2.0":
		return "licenses/apache-2.0", true
	case "BSD-3-Clause":
		return "licenses/bsd-3-clause", true
	case "MPL-2.0":
		return "licenses/mpl-2.0", true
	default:
		return "", false
	}
}

func sanitisePkgName(name string) string {
	base := filepath.Base(name)
	return strings.ReplaceAll(base, "-", "_")
}

func currentGoVersion() string {
	return goVersion
}

func currentYear() int {
	return time.Now().Year()
}
