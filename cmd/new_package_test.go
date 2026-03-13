package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestNewPackageCmd_DefaultFlags(t *testing.T) {
	dir := t.TempDir()
	projectDir := filepath.Join(dir, "testpkg")

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(origWd) })

	cmd := NewRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"new", "package", "testpkg"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
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
		"testpkg.go",
		"testpkg_test.go",
		"README.md",
		"AGENTS.md",
		"Makefile",
		filepath.Join(".github", "workflows", "ci.yml"),
		".scaff.json",
	}
	for _, f := range expectedFiles {
		full := filepath.Join(projectDir, f)
		if _, err := os.Stat(full); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist", f)
		}
	}
}

func TestNewPackageCmd_DisableFeatures(t *testing.T) {
	dir := t.TempDir()
	projectDir := filepath.Join(dir, "minimal")

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(origWd) })

	cmd := NewRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"new", "package", "minimal", "--agents=false", "--makefile=false", "--ci=false"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	shouldNotExist := []string{"AGENTS.md", "Makefile", filepath.Join(".github", "workflows", "ci.yml")}
	for _, f := range shouldNotExist {
		full := filepath.Join(projectDir, f)
		if _, err := os.Stat(full); err == nil {
			t.Errorf("file %s should not exist", f)
		}
	}

	shouldExist := []string{
		".editorconfig",
		".gitignore",
		".golangci.yml",
		"CHANGELOG.md",
		"CONTRIBUTING.md",
		"COPYRIGHT",
		"SECURITY.md",
	}
	for _, f := range shouldExist {
		full := filepath.Join(projectDir, f)
		if _, err := os.Stat(full); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist", f)
		}
	}
}

func TestNewPackageCmd_MissingArg(t *testing.T) {
	cmd := NewRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"new", "package"})

	if err := cmd.Execute(); err == nil {
		t.Error("expected error for missing argument")
	}
}

func TestNewPackageCmd_LicenseFlag(t *testing.T) {
	dir := t.TempDir()
	projectDir := filepath.Join(dir, "mitpkg")

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(origWd) })

	cmd := NewRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"new", "package", "mitpkg", "--license=MIT", "--agents=false", "--makefile=false", "--ci=false"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(projectDir, "LICENSE")); os.IsNotExist(err) {
		t.Fatal("expected LICENSE file to exist")
	}
	if _, err := os.Stat(filepath.Join(projectDir, "COPYRIGHT")); err == nil {
		t.Fatal("did not expect COPYRIGHT file to exist")
	}
}
