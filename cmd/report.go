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

// Static Mode

//Package cmd Handle report command line
package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/jsenon/worker-ops/config"
	"github.com/jsenon/worker-ops/internal/generate"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var before int8
var credfile string

//reportCmd launch static report
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Launch reporting on stdout",
	Long: `Launch a static report of Worker   
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

		StartReporter()
	},
}

//init initialize cobra command
func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.PersistentFlags().Int8Var(&before, "before", 12, "Specified before uptime in hour report will be generate")
	err := viper.BindPFlag("before", reportCmd.PersistentFlags().Lookup("before"))
	if err != nil {
		log.Error().Msgf("Error binding before value: %v", err.Error())
	}
	viper.SetDefault("before", 12)
	reportCmd.PersistentFlags().StringVar(&credfile, "credfile", "~/.aws/credentials", "Specify aws credential file")
	err = viper.BindPFlag("credfile", reportCmd.PersistentFlags().Lookup("credfile"))
	if err != nil {
		log.Error().Msgf("Error binding credfile value: %v", err.Error())
	}
	viper.SetDefault("credfile", "/.aws/credentials")
}

//StartReporter Start Static Report
func StartReporter() {
	ctx := context.Background()
	var stdoutmsg string
	vartime := viper.GetInt("before")

	remotejaegurl := viper.GetString("jaegerurl")
	if remotejaegurl != "" {
		log.Debug().Msg("Jaeger endpoint has been defined")
		// jaegerexporter.NewExporterCollector()
		// trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	}

	myresult, err := generate.Launch(ctx, vartime)
	if err != nil {
		log.Error().Msgf("Error Generate report: %v", err)
	}
	if len(myresult) > 0 {
		for _, n := range myresult {

			stdoutmsg = "Worker still runing since " + viper.GetString("before") + " hour(s) in " + n.Env + " environment: \n"

			for _, o := range n.Instances {
				stdoutmsg = stdoutmsg + " Name: " + o.Name + " DNS: " + o.Dnsname + " on region: " + o.Region + " Started since UTC: " + o.Launchtime.String() + "\n"
			}

		}

		fmt.Println(stdoutmsg)
	} else {
		fmt.Printf("No worker running since: %v hour(s)\n", viper.GetInt("before"))
	}
}
