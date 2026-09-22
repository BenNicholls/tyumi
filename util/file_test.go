package util

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestGetFileList(t *testing.T) {
	dir := t.TempDir()

	files := []string{"a.txt", "b.txt", "c.png"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(dir, f), []byte("data"), 0644); err != nil {
			t.Fatalf("failed to set up test file %s: %v", f, err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "subdir"), 0755); err != nil {
		t.Fatalf("failed to set up test subdir: %v", err)
	}

	t.Run("no extension filter", func(t *testing.T) {
		got, err := GetFileList(dir, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		slices.Sort(got)
		want := []string{"a.txt", "b.txt", "c.png"}
		if !slices.Equal(got, want) {
			t.Errorf("GetFileList(%q, \"\") = %v, want %v", dir, got, want)
		}
	})

	t.Run("with extension filter", func(t *testing.T) {
		got, err := GetFileList(dir, ".txt")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		slices.Sort(got)
		want := []string{"a.txt", "b.txt"}
		if !slices.Equal(got, want) {
			t.Errorf("GetFileList(%q, \".txt\") = %v, want %v", dir, got, want)
		}
	})

	t.Run("nonexistent directory", func(t *testing.T) {
		_, err := GetFileList(filepath.Join(dir, "does-not-exist"), "")
		if err == nil {
			t.Error("expected an error for a nonexistent directory")
		}
	})
}
