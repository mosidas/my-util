/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"util/password"

	"github.com/spf13/cobra"
)

var length int
var policy string

// passCmd represents the pass command
var passCmd = &cobra.Command{
	Use:   "pass",
	Short: "make a password",
	Long:  `make a password`,
	Run: func(cmd *cobra.Command, args []string) {
		// length >= 8
		if length < 8 {
			fmt.Println("length must be 8 or more")
			return
		}

		var policyInt int
		switch policy {
		case "all":
			policyInt = password.PolicyAllChars
		case "alphanum":
			policyInt = password.PolicyAlphaNum
		default:
			fmt.Println("invalid policy. use 'all' or 'alphanum'")
			return
		}

		password, err := password.MakePassword(length, policyInt)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(password)
	},
}

func init() {
	rootCmd.AddCommand(passCmd)

	passCmd.Flags().IntVarP(&length, "length", "l", 8, "password length")
	passCmd.Flags().StringVarP(&policy, "policy", "p", "all", "password policy (all or alphanum)")
}
