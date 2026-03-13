package filesystem_test

import (
	"os"
	"path/filepath"
	"testing"

	"scaff/internal/filesystem"
)

func TestWriteFile_CreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	w := filesystem.NewWriter(dir)

	if err := w.WriteFile(filepath.Join("a", "b", "c.txt"), "hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(dir, "a", "b", "c.txt"))
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}
	if string(content) != "hello" {
		t.Errorf("expected %q, got %q", "hello", string(content))
	}
}

func TestWriteFile_RefusesOverwrite(t *testing.T) {
	dir := t.TempDir()
	w := filesystem.NewWriter(dir)

	if err := w.WriteFile("file.txt", "first"); err != nil {
		t.Fatal(err)
	}
	if err := w.WriteFile("file.txt", "second"); err == nil {
		t.Error("expected error when overwriting, got nil")
	}
}

func TestEnsureDir_NonEmptyDir(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "x.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	w := filesystem.NewWriter(dir)
	if err := w.EnsureDir(); err == nil {
		t.Error("expected error for non-empty directory, got nil")
	}
}

func TestEnsureDir_CreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "newdir")
	w := filesystem.NewWriter(dir)

	if err := w.EnsureDir(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected directory")
	}
}
