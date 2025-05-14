/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nsclient",
	Short: "Test client for the ACORN Name Server",
	Args:  cobra.MinimumNArgs(1),
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

	var cfgFile string
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file. If not set, uses ./nsclient.yaml or $HOME/nsclient.yaml if exists.")
	rootCmd.PersistentFlags().String("addr", "", "gRPC server address. Defaults to localhost:8443 or localhost:8080 if --no-tls is set")
	rootCmd.PersistentFlags().String("ssl-cert", "", "SSL Certificate")
	rootCmd.PersistentFlags().Bool("no-tls", false, "Turn off tls")
	rootCmd.PersistentFlags().Bool("no-auth", false, "Turn off authentication")
	rootCmd.PersistentFlags().String("keycloak-url", "", "Keycloak server URL. For example: http://localhost:39507")
	rootCmd.PersistentFlags().String("user", "", "User")
	rootCmd.PersistentFlags().String("password", "", "Password")
	rootCmd.PersistentFlags().String("auth-client-id", "", "Client ID for authentication")
	rootCmd.PersistentFlags().String("auth-client-secret", "", "Client secret for authentication")
	rootCmd.PersistentFlags().String("json", "", "JSON input")
	rootCmd.PersistentFlags().Bool("json-template", false, "Create JSON template")
	rootCmd.PersistentFlags().Bool("verbose", false, "Print debug output")
	rootCmd.PersistentFlags().String("auth-realm", "", "Realm for authentication")

	// Bind flags to Viper
	if err := viper.BindPFlag("addr", rootCmd.PersistentFlags().Lookup("addr")); err != nil {
		log.Fatalf("Error binding flag 'addr': %v", err)
	}
	if err := viper.BindPFlag("ssl-cert", rootCmd.PersistentFlags().Lookup("ssl-cert")); err != nil {
		log.Fatalf("Error binding flag 'ssl-cert': %v", err)
	}
	if err := viper.BindPFlag("keycloak-url", rootCmd.PersistentFlags().Lookup("keycloak-url")); err != nil {
		log.Fatalf("Error binding flag 'auth-server': %v", err)
	}
	if err := viper.BindPFlag("auth.client-id", rootCmd.PersistentFlags().Lookup("auth-client-id")); err != nil {
		log.Fatalf("Error binding flag 'auth-client-id': %v", err)
	}
	if err := viper.BindPFlag("auth.realm", rootCmd.PersistentFlags().Lookup("auth-realm")); err != nil {
		log.Fatalf("Error binding flag 'auth-realm': %v", err)
	}
	if err := viper.BindPFlag("auth.client-secret", rootCmd.PersistentFlags().Lookup("auth-client-secret")); err != nil {
		log.Fatalf("Error binding flag 'auth-client-secret': %v", err)
	}
	if err := viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user")); err != nil {
		log.Fatalf("Error binding flag 'user': %v", err)
	}
	if err := viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password")); err != nil {
		log.Fatalf("Error binding flag 'password': %v", err)
	}

	// Configure Viper to read from a config file
	viper.SetConfigName("nsclient") // Name of the config file (without extension)
	viper.SetConfigType("yaml")     // Config file format
	viper.AddConfigPath(".")        // Look for the config file in the current directory
	viper.AddConfigPath("$HOME")    // Look for the config file in the user's home directory

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No configuration file found: %v", err)
	}
}
