/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
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
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {
	//	path := args[0]
	//	logrus.Debug("bla")
	//	doSomeStuff(path)
	//},
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.asyncApiCodeGen.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
