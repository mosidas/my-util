package zipfolder

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// zip に含めないエントリ
const (
	dsStoreFile = ".DS_Store"
	macosxDir   = "__MACOSX"
)

// srcDirの内容をzipPathに圧縮する
// zip内の各エントリはfilepath.Base(srcDir)をルートとする
// .DS_Storeファイルと__MACOSXディレクトリを除外する
// zipPathがsrcDir配下にある場合、そのzipファイル自身も除外する
// シンボリックリンク等の通常ファイル以外はスキップする
func ZipFolder(srcDir, zipPath string) error {
	srcDir, err := filepath.Abs(srcDir)
	if err != nil {
		return err
	}
	zipPath, err = filepath.Abs(zipPath)
	if err != nil {
		return err
	}

	info, err := os.Stat(srcDir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", srcDir)
	}

	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}

	w := zip.NewWriter(zipFile)
	parent := filepath.Dir(srcDir)

	walkErr := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() == macosxDir {
			return filepath.SkipDir
		}
		if !d.IsDir() && (d.Name() == dsStoreFile || path == zipPath) {
			return nil
		}

		fileInfo, err := d.Info()
		if err != nil {
			return err
		}
		if !d.IsDir() && !fileInfo.Mode().IsRegular() {
			return nil
		}

		relPath, err := filepath.Rel(parent, path)
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(relPath)

		if d.IsDir() {
			header.Name += "/"
			_, err = w.CreateHeader(header)
			return err
		}

		header.Method = zip.Deflate
		writer, err := w.CreateHeader(header)
		if err != nil {
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(writer, f)
		return err
	})

	if walkErr != nil {
		w.Close()
		zipFile.Close()
		os.Remove(zipPath)
		return walkErr
	}
	if err := w.Close(); err != nil {
		zipFile.Close()
		os.Remove(zipPath)
		return err
	}
	return zipFile.Close()
}
