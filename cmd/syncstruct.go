/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"gocanal/internal/canal"

	"github.com/spf13/cobra"
)

// syncstructCmd represents the syncstruct command
var syncstructCmd = &cobra.Command{
	Use:   "syncstruct",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		canal.SyncStruct(context.Background(), args...)
	},
}

func init() {
	rootCmd.AddCommand(syncstructCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncstructCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncstructCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
