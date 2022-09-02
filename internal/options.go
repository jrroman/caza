package internal

type Options struct {
	CloudEnabled bool   `long:"cloud" description:"Pull network data from a cloud provider"`
	MetricsPort  string `long:"metrics-port" description:"Port to serve up Prometheus metrics" default:":2112" env:"METRICS_ADDR"`
	Networks     string `long:"networks" description:"Comma separated list of name:cidr pairs; e.g. 'local:127.0.0.1/32,router:192.168.0.0/16'"`
	Region       string `long:"region" description:"Cloud region to use" env:"REGION"`
	VpcID        string `long:"vpc-id" description:"AWS vpc id" env:"AWS_VPC_ID"`
}
