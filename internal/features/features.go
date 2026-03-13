package features

import (
	"fmt"
	"path/filepath"
	"strings"

	"scaff/internal/filesystem"
	"scaff/internal/structure"
	"scaff/internal/templates"
)

type pack struct {
	component    string
	filename     string
	templateName string
}

var registry = map[string][]pack{
	"agents": {
		{component: "root", filename: "AGENTS.md", templateName: "agents/{project_type}.md"},
	},
	"makefile": {
		{component: "root", filename: "Makefile", templateName: "makefile"},
	},
	"ci": {
		{component: "ci", filename: "ci.yml", templateName: "ci.yml"},
	},
}

func Apply(w *filesystem.Writer, style structure.Style, data templates.Data, projectType string, featureNames []string) ([]string, error) {
	var written []string
	for _, name := range featureNames {
		packs, ok := registry[name]
		if !ok {
			return nil, fmt.Errorf("unknown feature pack: %q", name)
		}
		for _, p := range packs {
			tmplName := strings.ReplaceAll(p.templateName, "{project_type}", projectType)
			content, err := templates.Render(tmplName, data)
			if err != nil {
				return nil, fmt.Errorf("rendering feature %s/%s: %w", name, p.filename, err)
			}
			relPath := filepath.Join(style.Resolve(p.component), p.filename)
			if err := w.WriteFile(relPath, content); err != nil {
				return nil, fmt.Errorf("writing feature %s/%s: %w", name, p.filename, err)
			}
			written = append(written, relPath)
		}
	}
	return written, nil
}
