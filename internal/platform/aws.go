package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type AWSClient struct {
	client *ec2.Client
}

func NewAWSClient(region string) (*AWSClient, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	return &AWSClient{
		client: ec2.NewFromConfig(cfg),
	}, nil
}

// GetClusterInstance fetches instances based on the Project tag "nebula"
func (a *AWSClient) GetClusterInstance(projectName string) ([]ec2types.Instance, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("tag:Project"),
				Values: []string{projectName},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"running", "pending"},
			},
		},
	}

	result, err := a.client.DescribeInstances(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var instances []ec2types.Instance
	for _, reservation := range result.Reservations {
		instances = append(instances, reservation.Instances...)
	}
	return instances, nil
}

func (a *AWSClient) GetCurrentSpotPrice(instanceType types.InstanceType, az string) (float64, error) {
	input := &ec2.DescribeSpotPriceHistoryInput{
		InstanceTypes:       []ec2types.InstanceType{instanceType},
		AvailabilityZone:    aws.String(az),
		ProductDescriptions: []string{"Linux/UNIX"},
		StartTime:           aws.Time(time.Now().Add(-1 * time.Hour)),
		MaxResults:          aws.Int32(1),
	}

	result, err := a.client.DescribeSpotPriceHistory(context.TODO(), input)
	if err != nil || len(result.SpotPriceHistory) == 0 {
		return 0.0, fmt.Errorf("could not fetch spot price: %w", err)
	}

	var price float64
	fmt.Sscanf(*result.SpotPriceHistory[0].SpotPrice, "%f", &price)

	return price, nil
}

func (a *AWSClient) CheckForInterruptions(instanceID string) (bool, string, error) {
	input := &ec2.DescribeInstanceStatusInput{
		InstanceIds: []string{instanceID},
	}

	result, err := a.client.DescribeInstanceStatus(context.TODO(), input)
	if err != nil {
		return false, "", err
	}

	for _, status := range result.InstanceStatuses {
		for _, event := range status.Events {
			if event.Code == ec2types.EventCodeInstanceStop {
				return true, *event.Description, nil
			}
		}
	}

	return false, "", nil
}
