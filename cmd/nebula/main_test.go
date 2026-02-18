package main

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// Mock client
type mockPriceClient struct{}

func (m mockPriceClient) DescribeSpotPriceHistory(ctx context.Context, params *ec2.DescribeSpotPriceHistoryInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSpotPriceHistoryOutput, error) {
	return &ec2.DescribeSpotPriceHistoryOutput{
		SpotPriceHistory: []types.SpotPrice{
			{SpotPrice: aws.String("0.0500")},
		},
	}, nil
}

func TestGetSpotPrice(t *testing.T) {
	client := mockPriceClient{}
	price, err := getSpotPrice(context.TODO(), client, types.InstanceTypeT4gXlarge, "us-west-1a")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if price != 0.05 {
		t.Errorf("Expected 0.05, got %f", price)
	}
}
