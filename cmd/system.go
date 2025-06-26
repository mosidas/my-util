package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"util/system"
)

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "command for system",
	Long:  `command for system`,
}

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "システムイベントをリッスンする",
	Long: `電源イベントとセッションイベントをリアルタイムでリッスンします。

以下のイベントを検知します：
- スリープ/復帰
- 画面ロック/アンロック
- ログオン/ログオフ
- その他のシステムイベント`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("システムイベントリスナーを開始します...")
		fmt.Println("終了するには Ctrl+C を押してください")
		fmt.Println()

		if err := system.ListenSystemEvents(); err != nil {
			log.Fatalf("エラーが発生しました: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(systemCmd)
	systemCmd.AddCommand(listenCmd)
}
