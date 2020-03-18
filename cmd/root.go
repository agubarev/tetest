/*
Copyright © 2020 Andrei Gubarev <agubarev@protonmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/agubarev/tetest/internal/currency"
	"github.com/gocraft/dbr/v2"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"log"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var manager *currency.Manager

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tetest",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tetest.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// initialzing currency project
	initManager()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".tetest" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tetest")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initManager() {
	// basic validation
	feedURL := strings.TrimSpace(os.Getenv("FEED_URL"))
	if feedURL == "" {
		log.Fatal("env `FEED_URL` is not set")
	}

	// initializing main logger
	// NOTE: using a preset logger is sufficient for testing purposes
	l, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to initialize logger: %s", err)
	}

	//---------------------------------------------------------------------------
	// initializing mysql store
	//---------------------------------------------------------------------------
	l.Info("initialzing database connection")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	connection, err := dbr.Open("mysql", dsn, nil)
	if err != nil {
		log.Fatalf("failed to initialize mysql connection: %s", err)
	}

	l.Info("initializing default MySQL store")
	mysqlStore, err := currency.NewDefaultMySQLStore(connection)
	if err != nil {
		log.Fatalf("failed to initialize mysql backend store: %s", err)
	}

	//---------------------------------------------------------------------------
	// initialzing new currency manager
	//---------------------------------------------------------------------------
	l.Info("initializing currency manager")
	manager, err = currency.NewManager(mysqlStore, feedURL)
	if err != nil {
		log.Fatalf("failed to initialize currency manager: %s", err)
	}

	// assigning logger to the manager
	l.Info("configuring main logger")
	if err = manager.SetLogger(l); err != nil {
		log.Fatalf("failed to set main logger: %s", err)
	}
}
