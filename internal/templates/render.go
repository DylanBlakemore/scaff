package templates

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"strings"
	"text/template"
)

//go:embed files
var fsys embed.FS

var allTemplates *template.Template

func init() {
	allTemplates = template.New("")
	err := fs.WalkDir(fsys, "files", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}
		_, err = allTemplates.New(path).Parse(string(data))
		if err != nil {
			return fmt.Errorf("parsing %s: %w", path, err)
		}
		return nil
	})
	if err != nil {
		panic("parsing embedded templates: " + err.Error())
	}
}

type Data struct {
	Name       string
	ModulePath string
	GoVersion  string
	Year       int
}

func Render(name string, data Data) (string, error) {
	fullName := "files/" + name + ".tmpl"

	var buf bytes.Buffer
	if err := allTemplates.ExecuteTemplate(&buf, fullName, data); err != nil {
		return "", fmt.Errorf("executing template %s: %w", name, err)
	}
	return buf.String(), nil
}
