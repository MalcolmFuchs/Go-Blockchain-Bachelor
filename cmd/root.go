package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ega-blockchain",
	Short: "EGA Blockchain Application",
	Long:  `EGA Blockchain is a distributed application to manage patient records using blockchain technology.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
