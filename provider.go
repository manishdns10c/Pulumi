package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createProvider(ctx *pulumi.Context) (*aws.Provider, error) {
	// Configure the AWS provider
	provider, err := aws.NewProvider(ctx, "aws", &aws.ProviderArgs{
		Region:  pulumi.StringPtr("eu-central-1"), // specify the desired region
		Profile: pulumi.StringPtr("playground"),
	})
	if err != nil {
		return nil, err
	}
	return provider, nil
}
