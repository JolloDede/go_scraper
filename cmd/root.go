/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/JolloDede/go_scraper/src"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go_scraper",
	Short: "Scraping URLs and prints the 404 links.",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		verbose := cmd.Flag("verbose").Value.String() == "true"
		dic := src.HandleUrl(url, bool(verbose))

		for key, value := range dic {
			// if value >= 400 && value < 500 {
			// 	fmt.Println(key, " ", value)
			// }
			fmt.Println("Key:", key, "Value:", value)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var url string

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go_scraper.yaml)")
	rootCmd.PersistentFlags().StringVarP(&url, "url", "u", "http://localhost:80", "Root of the website you whant to check.")
	rootCmd.MarkFlagRequired("url")

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "When recursively Checking the Website.")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
