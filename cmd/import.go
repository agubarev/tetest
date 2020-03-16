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
	"log"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Imports currency feed",
	Run: func(cmd *cobra.Command, args []string) {
		feedURL := strings.TrimSpace(os.Getenv("TETEST_FEED_URL"))
		if feedURL == "" {
			log.Fatal("env `TETEST_FEED_URL` is not set")
		}

		// initializing feed parser
		p := gofeed.NewParser()

		f, err := p.ParseURL(feedURL)
		if err != nil {
			log.Fatalf("failed to parse feed [%s]: %s", feedURL, err)
		}

		for _, v := range f.Items {
			spew.Dump(strings.SplitN(strings.TrimSpace(v.Description), " ", 3))
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
