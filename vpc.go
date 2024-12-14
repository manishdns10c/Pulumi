package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createVPC(ctx *pulumi.Context, provider *aws.Provider) error {
	// Create a new VPC in the specified region
	idp_vpc, err := ec2.NewVpc(ctx, "idp", &ec2.VpcArgs{
		CidrBlock:          pulumi.String("10.0.0.0/16"),
		EnableDnsHostnames: pulumi.BoolPtr(true),
		Tags: pulumi.StringMap{
			"Name": pulumi.String("idp-vpc"),
		},
	}, pulumi.Provider(provider))

	if err != nil {
		return err
	}

	_, err_subnet_public_1a := ec2.NewSubnet(ctx, "subnet-public-1a", &ec2.SubnetArgs{
		VpcId:     idp_vpc.ID(),
		CidrBlock: pulumi.String("10.0.1.0/24"),
		Tags: pulumi.StringMap{
			"Name": pulumi.String("subnet-public-1a"),
		},
		MapPublicIpOnLaunch: pulumi.Bool(true),
		AvailabilityZone:    pulumi.String("eu-central-1a"),
	}, pulumi.Provider(provider))
	if err_subnet_public_1a != nil {
		return err_subnet_public_1a
	}

	_, err_subnet_public_1b := ec2.NewSubnet(ctx, "subnet-public-1b", &ec2.SubnetArgs{
		VpcId:            idp_vpc.ID(),
		CidrBlock:        pulumi.String("10.0.2.0/24"),
		AvailabilityZone: pulumi.String("eu-central-1b"),
		Tags: pulumi.StringMap{
			"Name": pulumi.String("subnet-public-1b"),
		},
		MapPublicIpOnLaunch: pulumi.Bool(true),
	}, pulumi.Provider(provider))
	if err_subnet_public_1b != nil {
		return err_subnet_public_1b
	}

	private_subnet_1a, err_subnet_private_1a := ec2.NewSubnet(ctx, "subnet-private-1a", &ec2.SubnetArgs{
		VpcId:            idp_vpc.ID(),
		CidrBlock:        pulumi.String("10.0.3.0/24"),
		AvailabilityZone: pulumi.String("eu-central-1a"),
		Tags: pulumi.StringMap{
			"Name": pulumi.String("subnet-private-1a"),
		},
	}, pulumi.Provider(provider))
	if err_subnet_private_1a != nil {
		return err_subnet_private_1a
	}

	private_subnet_1b, err_subnet_private_1b := ec2.NewSubnet(ctx, "subnet-private-1b", &ec2.SubnetArgs{
		VpcId:            idp_vpc.ID(),
		CidrBlock:        pulumi.String("10.0.4.0/24"),
		AvailabilityZone: pulumi.String("eu-central-1b"),
		Tags: pulumi.StringMap{
			"Name": pulumi.String("subnet-private-1b"),
		},
	}, pulumi.Provider(provider))
	if err_subnet_private_1b != nil {
		return err_subnet_private_1b
	}
	// create route table
	private_route_table, err_private_rt := ec2.NewRouteTable(ctx, "Private-RT", &ec2.RouteTableArgs{
		VpcId: idp_vpc.ID(),

		Tags: pulumi.StringMap{
			"Name": pulumi.String("Private-RT"),
		},
	}, pulumi.Provider(provider))
	if err_private_rt != nil {
		return err_private_rt
	}

	// Route table association
	_, err_rtb_association_1a := ec2.NewRouteTableAssociation(ctx, "private-rtb-1a", &ec2.RouteTableAssociationArgs{
		SubnetId:     private_subnet_1a.ID(),
		RouteTableId: private_route_table.ID(),
	}, pulumi.Provider(provider))
	if err_rtb_association_1a != nil {
		return err_rtb_association_1a
	}

	_, err_rtb_association_1b := ec2.NewRouteTableAssociation(ctx, "private-rtb-1b", &ec2.RouteTableAssociationArgs{
		SubnetId:     private_subnet_1b.ID(),
		RouteTableId: private_route_table.ID(),
	}, pulumi.Provider(provider))
	if err_rtb_association_1b != nil {
		return err_rtb_association_1b
	}
	// Internet gateway
	igw, err_igw := ec2.NewInternetGateway(ctx, "igw", &ec2.InternetGatewayArgs{
		VpcId: idp_vpc.ID(),
		Tags: pulumi.StringMap{
			"Name": pulumi.String("idp-igw"),
		},
	}, pulumi.Provider(provider))
	if err_igw != nil {
		return err_igw
	}

	// add igw to main rtb
	_, err_igw_rtb := ec2.NewRoute(ctx, "r", &ec2.RouteArgs{
		RouteTableId:         idp_vpc.MainRouteTableId,
		DestinationCidrBlock: pulumi.String("0.0.0.0/0"),
		GatewayId:            igw.ID(),
	}, pulumi.Provider(provider))
	if err_igw_rtb != nil {
		return err_igw_rtb
	}
	// elastic ip for nat
	eip, err_eip := ec2.NewEip(ctx, "eip", &ec2.EipArgs{}, pulumi.Provider(provider))
	if err_eip != nil {
		return err_eip
	}
	// Nat gateway
	nat, err_nat := ec2.NewNatGateway(ctx, "idp-nat", &ec2.NatGatewayArgs{
		AllocationId: eip.ID(),
		SubnetId:     private_subnet_1a.ID(),
		Tags: pulumi.StringMap{
			"Name": pulumi.String("gw NAT"),
		},
	}, pulumi.DependsOn([]pulumi.Resource{
		igw,
	}), pulumi.Provider(provider))
	if err_nat != nil {
		return err_nat
	}
	// Route from private subnet to nat
	_, err_nat_rtb := ec2.NewRoute(ctx, "nat-route", &ec2.RouteArgs{
		RouteTableId:         private_route_table.ID(),
		DestinationCidrBlock: pulumi.String("0.0.0.0/0"),
		NatGatewayId:         nat.ID(),
	}, pulumi.Provider(provider))
	if err_nat_rtb != nil {
		return err_nat_rtb
	}
	return nil
}
