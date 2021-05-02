/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"os"
	"scanboxie/pkg/scanboxie"
	"scanboxie/web"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server <barcodeDirMapFilepath>",
	Short: "",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		barcodeDirMapFilepath := args[0]

		scanboxie, err := scanboxie.NewScanboxie(barcodeDirMapFilepath, scanboxieConfig.CommandSets)
		if err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(1)
		}

		fmt.Printf("Wait for events on %s\n", scanboxieConfig.Eventpath)
		if scanboxieConfig.Eventpath != "" {
			go scanboxie.ListenAndProcessEvents(scanboxieConfig.Eventpath)
		}

		fmt.Printf("Listen on %s\n", scanboxieConfig.Bindaddress)
		webapp := web.NewApp(scanboxie, scanboxieConfig.ImageDir)
		webapp.Serve(scanboxieConfig.Bindaddress)

	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
