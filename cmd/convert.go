// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"

	"github.com/aperezg/linkstobook/converter"

	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert output into a epub file",
	Run: func(cmd *cobra.Command, args []string) {
		f, _ := cmd.Flags().GetString("format")
		outputDir, _ := cmd.Flags().GetString("output")
		webFiles, _ := cmd.Flags().GetStringSlice("web")

		c, err := converter.NewConverter(f,
			converter.WithOutputDir(outputDir),
			converter.WithWebFiles(webFiles),
		)
		if err != nil {
			log.Fatal(err)
		}
		if err = c.Convert(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringP("format", "f", "epub", "Output format file")
	convertCmd.Flags().StringP("output", "o", "", "Output dir where tar.gz will be saved")
	convertCmd.Flags().StringSlice("web", []string{}, "external html files to convert into epub (comma separated)")
}
