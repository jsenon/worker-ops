//Package prometheus launch calculation on worker node and export to prometheus
package prometheus

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	awspkg "github.com/jsenon/worker-ops/pkg/aws"
	pkgprometheus "github.com/jsenon/worker-ops/pkg/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	ini "github.com/vaughan0/go-ini"
)

//LauchRecord will start thread for worker calculation
func LauchRecord(ctx context.Context) {

	config := os.Getenv("HOME") + viper.GetString("credfile")
	if _, err := os.Stat(config); os.IsNotExist(err) {
		log.Error().Msgf("No config file found at: %v", config)
		os.Exit(1)
	}

	// Load ini file with credentials
	file, err := ini.LoadFile(config)
	if err != nil {
		log.Error().Msgf("Fail to read file: %v", err)
		os.Exit(1)
	}

	for key, values := range file {
		account := key
		key := values["aws_access_key_id"]
		pass := values["aws_secret_access_key"]
		creds := credentials.NewStaticCredentials(key, pass, "")
		region := "us-east-1"
		sess := session.Must(session.NewSession())

		svc := ec2.New(sess, &aws.Config{
			Credentials: creds,
			Region:      aws.String(region),
		})
		// Retrieve all region
		regions, err := svc.DescribeRegions(&ec2.DescribeRegionsInput{})
		if err != nil {
			log.Error().Msgf("Fail describe region: %v", err)
		}
		// Loop over all regions
		for _, region := range regions.Regions {
			// Get Instances info
			metricValue, err := awspkg.CountInstances(ctx, creds, account, *region.RegionName)
			if err != nil {
				log.Error().Msgf("Fail describe region: %v", err)
			}
			pkgprometheus.WorkerNumber.With(prometheus.Labels{"env": account, "region": *region.RegionName}).Set(float64(metricValue))

		}
	}
}
