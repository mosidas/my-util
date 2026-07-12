/*
Copyright © 2026 mosida sspydery@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"util/zipfolder"

	"github.com/spf13/cobra"
)

// zipCmd represents the zip command
var zipCmd = &cobra.Command{
	Use:   "zip [target directory]",
	Short: "zip each folder in the target directory",
	Long: `zip each folder in the target directory.
Each folder is compressed into <folder name>.zip in the target directory.
Hidden folders are skipped. .DS_Store files and __MACOSX directories are excluded.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetDir := args[0]
		entries, err := os.ReadDir(targetDir)
		if err != nil {
			fmt.Println(err)
			return
		}

		count := 0
		for _, entry := range entries {
			// 隠しフォルダをスキップ
			if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
				continue
			}

			folderName := entry.Name()
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
	},
}

func init() {
	rootCmd.AddCommand(zipCmd)
}
