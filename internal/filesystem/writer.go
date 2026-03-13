package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	dirPerm  fs.FileMode = 0o755
	filePerm fs.FileMode = 0o644
)

type Writer struct {
	BaseDir string
}

func NewWriter(baseDir string) *Writer {
	return &Writer{BaseDir: baseDir}
}

func (w *Writer) WriteFile(relPath, content string) error {
	full := filepath.Join(w.BaseDir, relPath)

	if _, err := os.Stat(full); err == nil {
		return fmt.Errorf("file already exists: %s", relPath)
	}

	dir := filepath.Dir(full)
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	if err := os.WriteFile(full, []byte(content), filePerm); err != nil {
		return fmt.Errorf("writing %s: %w", relPath, err)
	}
	return nil
}

func (w *Writer) EnsureDir() error {
	info, err := os.Stat(w.BaseDir)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("path exists and is not a directory: %s", w.BaseDir)
		}

		entries, err := os.ReadDir(w.BaseDir)
		if err != nil {
			return fmt.Errorf("reading directory %s: %w", w.BaseDir, err)
		}
		if len(entries) > 0 {
			return fmt.Errorf("directory is not empty: %s", w.BaseDir)
		}
		return nil
	}

	if os.IsNotExist(err) {
		return os.MkdirAll(w.BaseDir, dirPerm)
	}
	return fmt.Errorf("stat %s: %w", w.BaseDir, err)
}
