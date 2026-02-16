package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2DescribePriceAPI interface {
	DescribeSpotPriceHistory(ctx context.Context, params *ec2.DescribeSpotPriceHistoryInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSpotPriceHistoryOutput, error)
}

func getSpotPrice(ctx context.Context, client EC2DescribePriceAPI, instanceType types.InstanceType, az string) (float64, error) {
	result, err := client.DescribeSpotPriceHistory(ctx, &ec2.DescribeSpotPriceHistoryInput{
		InstanceTypes:       []types.InstanceType{instanceType},
		AvailabilityZone:    aws.String(az),
		ProductDescriptions: []string{"Linux/UNIX"},
		StartTime:           aws.Time(time.Now()),
	})
	if err != nil || len(result.SpotPriceHistory) == 0 {
		return 0, err
	}

	// SpotPrice returned as string.
	price, err := strconv.ParseFloat(*result.SpotPriceHistory[0].SpotPrice, 64)
	return price, err
}

func main() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-west-1"),
		config.WithSharedConfigProfile("nebula-homelab"),
	)
	if err != nil {
		log.Fatalf("Unable to load DSK config: %v", err)
	}

	client := ec2.NewFromConfig(cfg)

	// search for tagged resources "nebula".
	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{Name: aws.String("tag:Project"), Values: []string{"nebula"}},
			{Name: aws.String("instance-state-name"), Values: []string{"running"}},
		},
	}
	result, err := client.DescribeInstances(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	var ids []string
	var totalSessionCost float64

	for _, res := range result.Reservations {
		for _, ins := range res.Instances {
			ids = append(ids, *ins.InstanceId)

			currentPrice, _ := getSpotPrice(ctx, client, ins.InstanceType, *ins.Placement.AvailabilityZone)

			uptime := time.Since(*ins.LaunchTime)
			cost := uptime.Hours() * currentPrice
			totalSessionCost += cost

			fmt.Printf("Node: %s (%s)\n", *ins.InstanceId, ins.InstanceType)
			fmt.Printf("  └─ Live Price: $%.4f/hr\n", currentPrice)
			fmt.Printf("  └─ Uptime:     %.2f hours\n", uptime.Hours())
			fmt.Printf("  └─ Spent:      $%.4f\n\n", cost)
		}
	}

	if len(ids) == 0 {
		fmt.Println("No active nodes.")
		return
	}

	fmt.Printf("Total Session Cost: $%.2f\n", totalSessionCost)
	fmt.Print("Terminate all nodes? (y/n): ")
	var confirm string
	_, err = fmt.Scanln(&confirm)
	if err != nil {
		fmt.Printf("Confirm error: %v", err)
	}

	if confirm == "y" {
		_, err = client.TerminateInstances(ctx, &ec2.TerminateInstancesInput{InstanceIds: ids})
		if err != nil {
			fmt.Printf("Failed to send termination signals: %v", err)
		}

		fmt.Println("Termination signals sent.")
	}
}
