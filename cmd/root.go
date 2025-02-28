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

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go_scraper",
	Short: "Scraping URLs and prints the 404 links.",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		verbose := cmd.Flag("verbose").Value.String() == "true"
		uMap := src.HandleUrl(url, bool(verbose))

		for _, value := range uMap {
			if value.StatusCode >= 400 && value.StatusCode < 500 {
				fmt.Print(Red)
			}
			if value.StatusCode >= 200 && value.StatusCode < 300 {
				fmt.Print(Green)
			}
			fmt.Println("On ", value.Link.Url, "Link Text: ", value.Link.Content, "Value: ", value.StatusCode)
			fmt.Print(Reset)
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
