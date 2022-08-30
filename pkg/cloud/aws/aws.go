package awscloud

import (
	"net"

	"github.com/jrroman/caza/pkg/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// TODO use policy / role for credentials so this can work from local, an ec2 instance,
// or a kubernetes cluster.

// TODO let users specify their own tags to filter subnets with

type AwsCloudClient struct {
	ec2 ec2iface.EC2API
}

func New(cfg *config.Config) (*AwsCloudClient, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.Region),
	})
	if err != nil {
		return nil, err
	}
	return &AwsCloudClient{
		ec2: ec2.New(sess),
	}, nil
}

func (cc *AwsCloudClient) GetNetworks(cfg *config.Config) (map[string]*net.IPNet, error) {
	input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: aws.StringSlice([]string{cfg.VpcID}),
			},
		},
	}
	output, err := cc.ec2.DescribeSubnets(input)
	if err != nil {
		return nil, err
	}
	networks := make(map[string]*net.IPNet)
	for _, subnet := range output.Subnets {
		_, ipNet, err := net.ParseCIDR(*subnet.CidrBlock)
		if err != nil {
			return nil, err
		}
		networks[*subnet.AvailabilityZone] = ipNet
	}
	return networks, nil
}
