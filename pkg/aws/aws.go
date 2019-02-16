//Package aws will Get token, initialize session, Retrieve Instances
package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
)

// Instance struct represent instance with its name, dns name, launch time and region
type Instance struct {
	Name       string    `json:"name"`
	Dnsname    string    `json:"dnsname"`
	Launchtime time.Time `json:"launchtime"`
	Region     string    `json:"region"`
}

// RetrieveInstances print all instances associated with credential on all region, some filters is applied staticaly
func RetrieveInstances(ctx context.Context, sp opentracing.Span, creds *credentials.Credentials, account string, region string, vartime int) (instances []Instance, err error) {
	span := opentracing.StartSpan(
		"(*worker-ops).RetrieveInstances",
		opentracing.ChildOf(sp.Context()))
	defer span.Finish()
	sess := session.Must(session.NewSession())

	// instances := []Instance{}

	svc := ec2.New(sess, &aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	})
	// Filter parameters
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Type"),
				Values: []*string{aws.String("worker")},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running"), aws.String("pending")},
			},
		},
	}
	resp, err := svc.DescribeInstances(params)
	if err != nil {
		log.Error().Msgf("Fail describe instances: %v", err)
		return nil, err
	}

	// Loop over instances
	for idx := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			name := "None"
			// Loop over tags
			for _, keys := range inst.Tags {
				if *keys.Key == "Name" {
					name = *keys.Value
				}
			}
			// calcutate instances that run before hour specified by command line (default 12)
			delay := time.Duration(vartime) * time.Hour
			newtime := time.Now().Add(-delay)
			if inst.LaunchTime.Before(newtime) {
				m := Instance{name, *inst.PrivateDnsName, *inst.LaunchTime, region}
				instances = append(instances, m)
				log.Debug().Msgf("Info: %v %v %v %v", name, *inst.PrivateDnsName, *inst.LaunchTime, region)
			}
		}
	}

	return instances, nil
}

// CountInstances calculate number of instances associated with credential on all region, some filters is applied staticaly
func CountInstances(ctx context.Context, sp opentracing.Span, creds *credentials.Credentials, account string, region string) (nbr int, err error) {
	span := opentracing.StartSpan(
		"(*worker-ops).CountInstances",
		opentracing.ChildOf(sp.Context()))
	defer span.Finish()
	sess := session.Must(session.NewSession())
	svc := ec2.New(sess, &aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	})
	// Filter parameters
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Type"),
				Values: []*string{aws.String("worker")},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running"), aws.String("pending")},
			},
		},
	}
	resp, err := svc.DescribeInstances(params)
	if err != nil {
		log.Error().Msgf("Fail describe instances: %v", err)
		return 0, err
	}
	nbr = len(resp.Reservations)
	return nbr, nil
}
