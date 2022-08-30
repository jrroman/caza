package internal

type Options struct {
	CloudEnabled bool   `long:"cloud" description:"Pull network data from a cloud provider"`
	MetricsPort  string `long:"metrics-port" description:"Port to serve up Prometheus metrics" default:":2112" env:"METRICS_ADDR"`
	Networks     string `long:"networks" description:"Comma separated string of CIDR's" env:"NETWORKS"`
	Region       string `long:"region" description:"Cloud region to use" default:"us-east-1" env:"REGION"`
	VpcID        string `long:"vpc-id" description:"AWS vpc id" env:"AWS_VPC_ID"`
}
