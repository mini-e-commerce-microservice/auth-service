package main

import "github.com/spf13/cobra"

func main() {
	var rootCmd = &cobra.Command{}

	rootCmd.AddCommand(consumerCDCUserData)
	rootCmd.AddCommand(restApi)
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
