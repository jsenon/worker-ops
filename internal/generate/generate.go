//Package generate will generate report
package generate

import (
	"context"
	"os"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	ini "github.com/vaughan0/go-ini"

	awspkg "github.com/jsenon/worker-ops/pkg/aws"
)

// FullInstances represent instances link to a profile in the credential file
// swagger:model FullInstances
type FullInstances struct {
	// in: body
	Env       string            `json:"env"`
	Instances []awspkg.Instance `json:"instances"`
}

//Launch launch static generation of worker state
func Launch(ctx context.Context, sp opentracing.Span, vartime int) ([]FullInstances, error) {
	span := opentracing.StartSpan(
		"(*worker-ops).Launch",
		opentracing.ChildOf(sp.Context()))
	defer span.Finish()

	log.Debug().Msgf("Starting reporting on all regions with all your credentials")

	myfullinstances := []FullInstances{}

	config := os.Getenv("HOME") + viper.GetString("credfile")
	if _, err := os.Stat(config); os.IsNotExist(err) {
		log.Error().Msgf("No config file found at: %v", config)
		return nil, err
	}

	// Load ini file with credentials
	file, err := ini.LoadFile(config)
	if err != nil {
		log.Error().Msgf("Fail to read file: %v", err)
		return nil, err
	}

	// Loop over credential ini file
	for key, values := range file {

		account := key
		key := values["aws_access_key_id"]
		pass := values["aws_secret_access_key"]
		creds := credentials.NewStaticCredentials(key, pass, "")
		region := "us-east-1"
		sess := session.Must(session.NewSession())
		log.Debug().Msgf("Retrieve instance for environment: %v", account)

		svc := ec2.New(sess, &aws.Config{
			Credentials: creds,
			Region:      aws.String(region),
		})
		// Retrieve all region
		regions, err := svc.DescribeRegions(&ec2.DescribeRegionsInput{})
		if err != nil {
			log.Error().Msgf("Fail describe region: %v", err)
			return nil, err
		}
		// Loop over all regions
		for _, region := range regions.Regions {

			// Get Instances info
			instances, err := awspkg.RetrieveInstances(ctx, span, creds, account, *region.RegionName, vartime)
			if err != nil {
				log.Error().Msgf("Fail to get instance: %v", err)
				return nil, err
			}
			if len(instances) > 0 {
				m := FullInstances{account, instances}
				myfullinstances = append(myfullinstances, m)
			}
		}
	}
	if len(myfullinstances) > 0 {
		log.Debug().Msgf("Instances: %v ", myfullinstances)
	} else {
		log.Debug().Msgf("No instances found")
	}
	return myfullinstances, nil
}
