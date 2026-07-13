/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "util",
	Short: "utility tool",
	Long:  `utility tool`,
}

func Execute() {
	rootCmd.SetHelpCommand(&cobra.Command{Use: "no-help", Hidden: true})
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
