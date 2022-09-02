package awscloud

import (
	"net"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

var (
	mockDescribeSubnetsOutput = &ec2.DescribeSubnetsOutput{
		Subnets: []*ec2.Subnet{
			{
				AvailabilityZone: aws.String("us-east-1a"),
				CidrBlock:        aws.String("10.0.0.1/16"),
			},
		},
	}
)

type mockDescribeSubnetsAPI func(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error)

type mockEC2Client struct {
	ec2iface.EC2API
	response *ec2.DescribeSubnetsOutput
}

func (c mockEC2Client) DescribeSubnets(input *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	return c.response, nil
}

// TODO create test helpers
func createIPNetHelper(cidrBlock string) *net.IPNet {
	_, ipNet, _ := net.ParseCIDR(cidrBlock)
	return ipNet
}

func TestGetNetworks(t *testing.T) {
	cases := []struct {
		name      string
		vpcID     string
		expect    map[string]*net.IPNet
		wantError bool
	}{
		{
			name:  "valid describe subnets",
			vpcID: "abc123",
			expect: map[string]*net.IPNet{
				"us-east-1a": createIPNetHelper("10.0.0.1/16"),
			},
			wantError: false,
		},
	}
	mockEC2API := mockEC2Client{response: mockDescribeSubnetsOutput}
	mockClient := &AwsCloudClient{ec2: mockEC2API}
	for _, tc := range cases {
		got, err := mockClient.GetNetworks("abc123")
		if tc.wantError && err == nil {
			t.Errorf("testcase: %s expected an error", tc.name)
			continue
		}
		if !tc.wantError && err != nil {
			t.Errorf("testcase: %s did not expect error", tc.name)
			continue
		}
		if !reflect.DeepEqual(got, tc.expect) {
			t.Errorf("testcase: %s; got: %v, want: %v", tc.name, got, tc.expect)
		}
	}
}
