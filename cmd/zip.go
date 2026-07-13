/*
Copyright © 2026 mosida sspydery@gmail.com
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"util/zipfolder"

	"github.com/spf13/cobra"
)

var skipConfirm bool

// zipCmd represents the zip command
var zipCmd = &cobra.Command{
	Use:   "zip [target directory]",
	Short: "zip folders",
	Long: `zip each folder in the target directory.
Each folder is compressed into <folder name>.zip in the target directory.
If no target directory is given, the current directory is compressed
into <current directory name>.zip in the current directory.
Hidden folders are skipped. .DS_Store files and __MACOSX directories are excluded.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			zipCurrentDir()
			return
		}
		zipEachFolder(args[0])
	},
}

// カレントディレクトリを<カレントディレクトリ名>.zipに圧縮する
func zipCurrentDir() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	zipName := filepath.Base(cwd) + ".zip"

	fmt.Printf("target: %s -> %s\n", cwd, zipName)
	if !confirm() {
		fmt.Println("canceled")
		return
	}

	if err := zipfolder.ZipFolder(cwd, filepath.Join(cwd, zipName)); err != nil {
		fmt.Printf("error: failed to compress: %v\n", err)
		return
	}
	fmt.Printf("done: created %s\n", zipName)
}

// targetDir直下の各フォルダを<フォルダ名>.zipに圧縮する
func zipEachFolder(targetDir string) {
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	var folders []string
	for _, entry := range entries {
		// 隠しフォルダをスキップ
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		folders = append(folders, entry.Name())
	}
	if len(folders) == 0 {
		fmt.Println("no folders to compress")
		return
	}

	fmt.Printf("target directory: %s\n", targetDir)
	fmt.Printf("folders to compress: %s\n", strings.Join(folders, ", "))
	if !confirm() {
		fmt.Println("canceled")
		return
	}

	count := 0
	for _, folderName := range folders {
		zipName := folderName + ".zip"
		fmt.Printf("compressing: %s -> %s\n", folderName, zipName)

		err := zipfolder.ZipFolder(filepath.Join(targetDir, folderName), filepath.Join(targetDir, zipName))
		if err != nil {
			fmt.Printf("error: failed to compress %s: %v\n", folderName, err)
			continue
		}
		count++
	}

	fmt.Printf("done: compressed %d folder(s)\n", count)
}

// 実行してよいかy/nで確認する
// -yオプション指定時は確認をスキップする
func confirm() bool {
	if skipConfirm {
		return true
	}
	fmt.Print("proceed? (y/n): ")
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return false
	}
	answer := strings.TrimSpace(strings.ToLower(line))
	return answer == "y" || answer == "yes"
}

func init() {
	rootCmd.AddCommand(zipCmd)

	zipCmd.Flags().BoolVarP(&skipConfirm, "yes", "y", false, "skip confirmation")
}
