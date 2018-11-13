package cmd

import (
	"fmt"

	"github.com/govenue/configurator"
	"github.com/govenue/goman"
	"github.com/govenue/homedir"
)

var (
	// Used for flags.
	cfgFile, userLicense string

	rootCmd = &goman.Command{
		Use:   "goman",
		Short: "A generator for goman based Applications",
		Long: `goman is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a goman application.`,
	}
)

// Execute executes the root command.
func Execute() {
	rootCmd.Execute()
}

func init() {
	goman.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goman.yaml)")
	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	rootCmd.PersistentFlags().Bool("configurator", true, "use configurator for configuration")
	configurator.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	configurator.BindPFlag("useConfigurator", rootCmd.PersistentFlags().Lookup("configurator"))
	configurator.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	configurator.SetDefault("license", "apache")

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(initCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		configurator.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		// Search config in home directory with name ".goman" (without extension).
		configurator.AddConfigPath(home)
		configurator.SetConfigName(".goman")
	}

	configurator.AutomaticEnv()

	if err := configurator.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", configurator.ConfigFileUsed())
	}
}
