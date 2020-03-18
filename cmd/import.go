/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/agubarev/tetest/internal/currency"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/v2"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Imports currency feed",
	Run: func(cmd *cobra.Command, args []string) {
		// basic validation
		feedURL := strings.TrimSpace(os.Getenv("FEED_URL"))
		if feedURL == "" {
			log.Fatal("env `FEED_URL` is not set")
		}

		//---------------------------------------------------------------------------
		// initializing mysql store
		//---------------------------------------------------------------------------
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

		mysqlStore, err := currency.NewDefaultMySQLStore(connection)
		if err != nil {
			log.Fatalf("failed to initialize mysql backend store: %s", err)
		}

		//---------------------------------------------------------------------------
		// initialzing new currency manager
		//---------------------------------------------------------------------------
		m, err := currency.NewManager(mysqlStore, feedURL)
		if err != nil {
			log.Fatalf("failed to initialize currency manager: %s", err)
		}

		// initializing main logger
		l, err := zap.NewProduction()
		if err != nil {
			log.Fatalf("failed to initialize logger: %s", err)
		}

		// assigning logger to the manager
		if err = m.SetLogger(l); err != nil {
			log.Fatalf("failed to set main logger: %s", err)
		}

		// importing currencies from remote source
		if err = m.Import(context.Background()); err != nil {
			log.Fatalf("failed to import currency: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
