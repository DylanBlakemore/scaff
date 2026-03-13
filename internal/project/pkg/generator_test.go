package pkg_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"scaff/internal/generator"
	"scaff/internal/metadata"
	"scaff/internal/project/pkg"
)

func TestGenerate_AllFeaturesEnabled(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "mylib")

	gen := pkg.New()
	result, err := gen.Generate(generator.Request{
		Name:       "mylib",
		ModulePath: "github.com/user/mylib",
		Style:      "default",
		Features:   []string{"agents", "makefile", "ci"},
		OutputDir:  projectDir,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedFiles := []string{
		".editorconfig",
		".gitignore",
		".golangci.yml",
		"CHANGELOG.md",
		"CONTRIBUTING.md",
		"COPYRIGHT",
		"SECURITY.md",
		"go.mod",
		"mylib.go",
		"mylib_test.go",
		"README.md",
		"AGENTS.md",
		"Makefile",
		filepath.Join(".github", "workflows", "ci.yml"),
		".scaff.json",
	}

	for _, f := range expectedFiles {
		if !slices.Contains(result.Files, f) {
			t.Errorf("expected file %q in result, got %v", f, result.Files)
		}
		full := filepath.Join(projectDir, f)
		if _, err := os.Stat(full); os.IsNotExist(err) {
			t.Errorf("expected file to exist on disk: %s", full)
		}
	}
}

func TestGenerate_NoFeatures(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "bare")

	gen := pkg.New()
	result, err := gen.Generate(generator.Request{
		Name:       "bare",
		ModulePath: "bare",
		Style:      "default",
		Features:   []string{},
		OutputDir:  projectDir,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	shouldExist := []string{
		".editorconfig",
		".gitignore",
		".golangci.yml",
		"CHANGELOG.md",
		"CONTRIBUTING.md",
		"COPYRIGHT",
		"SECURITY.md",
		"go.mod",
		"bare.go",
		"bare_test.go",
		"README.md",
		".scaff.json",
	}
	shouldNotExist := []string{"AGENTS.md", "Makefile", filepath.Join(".github", "workflows", "ci.yml")}

	for _, f := range shouldExist {
		if !slices.Contains(result.Files, f) {
			t.Errorf("expected file %q in result", f)
		}
		full := filepath.Join(projectDir, f)
		if _, err := os.Stat(full); os.IsNotExist(err) {
			t.Errorf("expected file on disk: %s", full)
		}
	}

	for _, f := range shouldNotExist {
		if slices.Contains(result.Files, f) {
			t.Errorf("did not expect file %q in result", f)
		}
		full := filepath.Join(projectDir, f)
		if _, err := os.Stat(full); err == nil {
			t.Errorf("file should not exist on disk: %s", full)
		}
	}
}

func TestGenerate_MetadataContents(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "metacheck")

	gen := pkg.New()
	_, err := gen.Generate(generator.Request{
		Name:       "metacheck",
		ModulePath: "metacheck",
		Style:      "default",
		Features:   []string{"makefile"},
		OutputDir:  projectDir,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw, err := os.ReadFile(filepath.Join(projectDir, metadata.FileName))
	if err != nil {
		t.Fatalf("reading metadata: %v", err)
	}

	var meta metadata.Project
	if err := json.Unmarshal(raw, &meta); err != nil {
		t.Fatalf("unmarshalling metadata: %v", err)
	}

	if meta.ProjectType != "package" {
		t.Errorf("expected project_type %q, got %q", "package", meta.ProjectType)
	}
	if meta.Style != "default" {
		t.Errorf("expected style %q, got %q", "default", meta.Style)
	}
	if !slices.Contains(meta.Features, "makefile") {
		t.Errorf("expected features to contain %q, got %v", "makefile", meta.Features)
	}
	if slices.Contains(meta.Features, "ci") {
		t.Errorf("features should not contain %q", "ci")
	}
}

func TestGenerate_LicenseSelected_WritesLICENSE(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "licensed")

	gen := pkg.New()
	result, err := gen.Generate(generator.Request{
		Name:       "licensed",
		ModulePath: "github.com/user/licensed",
		Style:      "default",
		Features:   []string{},
		OutputDir:  projectDir,
		License:    "MIT",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !slices.Contains(result.Files, "LICENSE") {
		t.Fatalf("expected LICENSE in result files, got %v", result.Files)
	}
	if slices.Contains(result.Files, "COPYRIGHT") {
		t.Fatalf("did not expect COPYRIGHT when LICENSE is generated")
	}

	if _, err := os.Stat(filepath.Join(projectDir, "LICENSE")); os.IsNotExist(err) {
		t.Fatal("expected LICENSE file to exist")
	}
	if _, err := os.Stat(filepath.Join(projectDir, "COPYRIGHT")); err == nil {
		t.Fatal("did not expect COPYRIGHT file to exist")
	}
}

func TestGenerate_HyphenatedName(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "my-lib")

	gen := pkg.New()
	result, err := gen.Generate(generator.Request{
		Name:       "my-lib",
		ModulePath: "github.com/user/my-lib",
		Style:      "default",
		Features:   []string{},
		OutputDir:  projectDir,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !slices.Contains(result.Files, "my_lib.go") {
		t.Errorf("expected sanitised filename my_lib.go, got %v", result.Files)
	}

	content, err := os.ReadFile(filepath.Join(projectDir, "my_lib.go"))
	if err != nil {
		t.Fatalf("reading source: %v", err)
	}
	if got := string(content); got != "package my_lib\n" {
		t.Errorf("unexpected package declaration: %q", got)
	}
}

func TestGenerate_ValidationErrors(t *testing.T) {
	gen := pkg.New()

	tests := []struct {
		name string
		req  generator.Request
	}{
		{"empty name", generator.Request{Name: "", ModulePath: "mod", Style: "default"}},
		{"empty module", generator.Request{Name: "x", ModulePath: "", Style: "default"}},
		{"bad style", generator.Request{Name: "x", ModulePath: "x", Style: "unknown"}},
		{"bad feature", generator.Request{Name: "x", ModulePath: "x", Style: "default", Features: []string{"nope"}}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.req.OutputDir = filepath.Join(t.TempDir(), "out")
			_, err := gen.Generate(tc.req)
			if err == nil {
				t.Error("expected validation error, got nil")
			}
		})
	}
}

func TestGenerate_RefusesNonEmptyDir(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "existing")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, "something.txt"), []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}

	gen := pkg.New()
	_, err := gen.Generate(generator.Request{
		Name:       "existing",
		ModulePath: "existing",
		Style:      "default",
		OutputDir:  projectDir,
	})
	if err == nil {
		t.Error("expected error for non-empty directory, got nil")
	}
}
