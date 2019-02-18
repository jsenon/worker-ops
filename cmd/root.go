// Copyright © 2019 Delair <julien.senon@delair.aero>
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

//Package cmd Handle root command line
package cmd

import (
	"fmt"
	"os"

	"github.com/jsenon/worker-ops/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var loglevel bool
var jaegerurl string
var version bool

//rootCmd is the root of the command line
var rootCmd = &cobra.Command{
	Use:   "worker-ops",
	Short: "Utility to report worker state",
	Long: `An Utility to launch reporter for worker node, or launch a web server API
to export usage to prometheus, and send report to slack`,
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			fmt.Printf("Version: %v, build from: %v, on: %v\n", config.Version, config.GitCommit, config.BuildDate)
		}
	},
}

//Execute will launch command line
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//init initialize cobra root command
func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVar(&loglevel, "debug", false, "Set log level to Debug") //nolint
	rootCmd.PersistentFlags().StringVar(&jaegerurl, "jaegerurl", "", "Set jaegger collector endpoint")
	rootCmd.PersistentFlags().BoolVar(&version, "version", false, "version")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	err := viper.BindPFlag("jaegerurl", rootCmd.PersistentFlags().Lookup("jaegerurl"))
	if err != nil {
		log.Error().Msgf("Error binding jaegerurl value: %v", err.Error())
	}
	viper.SetDefault("jaegerurl", "")

}

//initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".worker-ops")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
