package cmd

import (
	"github.com/c0olix/asyncApiCodeGen/logging"
	"os"

	"github.com/spf13/cobra"
)

var logger = logging.NewLogger()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "asyncApiCodeGen",
	Short: "Used to generate code for given async api spec",
	Long:  `This CLI-Tool is used to generate code for given async api spec`,
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
	rootCmd.PersistentFlags().StringP("input", "i", "", "Where is the source spec located?")
}
