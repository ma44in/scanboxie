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

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

type config struct {
	CommandSets *scanboxie.CommandSets
}

var cfgFile string
var conf *config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scanboxie <EventPath> <Barcode Music Dir Map Json File >",
	Short: "Control Music with Barcode Scanner",
	Long:  ``,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		eventPath := args[0]
		barcodeDirMapFilepath := args[1]

		watchForChanges, _ := cmd.Flags().GetBool("watch")

		scanboxie, err := scanboxie.NewScanboxie(barcodeDirMapFilepath, conf.CommandSets, watchForChanges)
		if err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(1)
		}

		fmt.Println("Listen and Process Events ...")
		err = scanboxie.ListenAndProcessEvents(eventPath)
		if err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(1)
		}
	},
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.scanboxie.yaml)")
	rootCmd.PersistentFlags().BoolP("watch", "w", false, "Watch for changes in barcodemap file")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AddConfigPath("/etc")
	viper.AddConfigPath(".")
	viper.SetConfigName("scanboxie")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	viper.AutomaticEnv() // read in environment variables that match

	conf = &config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}
}
