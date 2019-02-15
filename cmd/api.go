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

// API Mode

//Package cmd Handle server command line
package cmd

import (
	"context"
	"os"
	"runtime"

	"github.com/jsenon/worker-ops/config"
	"github.com/jsenon/worker-ops/pkg/rest"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//apiCmd will launch API server

var apibefore int8

var apiCmd = &cobra.Command{
	Use:   "server",
	Short: "Launch Worker Ops Server",
	Long: `Launch Worker Ops API Server 
           which manage worker status generate promtheus metrics and slack event
           `,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger = log.With().Str("Service", config.Service).Logger()
		log.Logger = log.With().Str("Version", config.Version).Logger()

		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		if loglevel {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
			err := os.Setenv("LOGLEVEL", "debug")
			if err != nil {
				log.Error().Msgf("Error %s", err.Error())
				runtime.Goexit()
			}
		}
		log.Debug().Msg("Log level set to Debug")
		log.Debug().Msgf("Jaeger Remote URL %s", viper.GetString("jaegerurl"))

		StartAPI()
	},
}

//init initialize cobra command
func init() {

	rootCmd.AddCommand(apiCmd)
	apiCmd.PersistentFlags().Int8Var(&apibefore, "apibefore", 1, "Specified before uptime in hour slack report will be generate")
	err := viper.BindPFlag("apibefore", apiCmd.PersistentFlags().Lookup("apibefore"))
	if err != nil {
		log.Error().Msgf("Error binding before value: %v", err.Error())
	}
	viper.SetDefault("apibefore", 1)
}

//StartAPI Start the server
func StartAPI() {
	ctx := context.Background()
	remotejaegurl := viper.GetString("jaegerurl")
	if remotejaegurl != "" {
		log.Debug().Msg("Jaeger endpoint has been defined")
		// jaegerexporter.NewExporterCollector()
		// trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	}

	rest.ServeRest(ctx)
}
