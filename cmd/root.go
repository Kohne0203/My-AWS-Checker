/*
Copyright © 2026 ko.watanabe <lsbean62@gmail.com
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "my-aws-checker",
	Short: "A CLI tool to audit AWS resource security configurations",
	Long: `my-aws-checker is a command-line tool for auditing AWS resources.
It scans your AWS environment and checks for security risks such as
misconfigured settings. Currently supports S3 bucket public access auditing.

Example usage:
  my-aws-checker s3              # Check S3 bucket security configurations
  my-aws-checker s3 --help       # Show detailed help for S3 command

Authentication:
Uses AWS credential chain (environment variables, ~/.aws/credentials, or IAM roles)`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.my-aws-checker.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
