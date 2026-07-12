package zipfolder

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestZipFolder(t *testing.T) {
	tmpDir := t.TempDir()

	// テスト用のフォルダ構成を作成
	srcDir := filepath.Join(tmpDir, "target")
	mustMkdir(t, filepath.Join(srcDir, "sub"))
	mustMkdir(t, filepath.Join(srcDir, "empty"))
	mustMkdir(t, filepath.Join(srcDir, "__MACOSX"))
	mustWrite(t, filepath.Join(srcDir, "a.txt"), "hello")
	mustWrite(t, filepath.Join(srcDir, "sub", "b.txt"), "world")
	mustWrite(t, filepath.Join(srcDir, ".DS_Store"), "junk")
	mustWrite(t, filepath.Join(srcDir, "sub", ".DS_Store"), "junk")
	mustWrite(t, filepath.Join(srcDir, "__MACOSX", "c.txt"), "junk")

	zipPath := filepath.Join(tmpDir, "target.zip")
	if err := ZipFolder(srcDir, zipPath); err != nil {
		t.Fatalf("ZipFolder failed: %v", err)
	}

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatalf("failed to open zip: %v", err)
	}
	defer r.Close()

	got := map[string]bool{}
	for _, f := range r.File {
		got[f.Name] = true
	}

	wants := []string{"target/", "target/a.txt", "target/sub/", "target/sub/b.txt", "target/empty/"}
	for _, want := range wants {
		if !got[want] {
			t.Errorf("entry %q not found in zip: %v", want, got)
		}
	}

	excludes := []string{"target/.DS_Store", "target/sub/.DS_Store", "target/__MACOSX/", "target/__MACOSX/c.txt"}
	for _, exclude := range excludes {
		if got[exclude] {
			t.Errorf("entry %q should be excluded from zip", exclude)
		}
	}
}

func TestZipFolderNotDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "file.txt")
	mustWrite(t, filePath, "hello")

	if err := ZipFolder(filePath, filepath.Join(tmpDir, "file.zip")); err == nil {
		t.Error("ZipFolder should fail for a non-directory")
	}
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
