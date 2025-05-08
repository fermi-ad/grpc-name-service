/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nsclient",
	Short: "Test client for the ACORN Name Server",
	Args:  cobra.MinimumNArgs(1),
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nsclient.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().String("addr", "", "gRPC server address. Defaults to localhost:8443 or localhost:8080 if --no-tls is set")
	rootCmd.PersistentFlags().String("ssl-cert", "../nameserver/src/main/resources/ssl/server.crt", "SSL Certificate")
	rootCmd.PersistentFlags().Bool("no-tls", false, "Turn off tls")
	rootCmd.PersistentFlags().Bool("no-authnz", false, "Turn off authentication and authorization")
	rootCmd.PersistentFlags().String("keycloak-url", "", "Keycloak URL. For example: http://localhost:39507")
	rootCmd.PersistentFlags().String("user", "alice", "User")
	rootCmd.PersistentFlags().String("password", "alice", "Password")
	rootCmd.PersistentFlags().String("json", "", "JSON input")
	rootCmd.PersistentFlags().Bool("json-template", false, "Create JSON template")
	rootCmd.PersistentFlags().Bool("verbose", false, "Print debug output")
}
