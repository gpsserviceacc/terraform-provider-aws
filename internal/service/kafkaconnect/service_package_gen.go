// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package kafkaconnect

import (
	"context"

	aws_sdkv1 "github.com/aws/aws-sdk-go/aws"
	session_sdkv1 "github.com/aws/aws-sdk-go/aws/session"
	kafkaconnect_sdkv1 "github.com/aws/aws-sdk-go/service/kafkaconnect"
	"terraform-provider-awsgps/internal/conns"
	"terraform-provider-awsgps/internal/types"
	"terraform-provider-awsgps/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  DataSourceConnector,
			TypeName: "aws_mskconnect_connector",
		},
		{
			Factory:  DataSourceCustomPlugin,
			TypeName: "aws_mskconnect_custom_plugin",
		},
		{
			Factory:  DataSourceWorkerConfiguration,
			TypeName: "aws_mskconnect_worker_configuration",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceConnector,
			TypeName: "aws_mskconnect_connector",
		},
		{
			Factory:  ResourceCustomPlugin,
			TypeName: "aws_mskconnect_custom_plugin",
		},
		{
			Factory:  ResourceWorkerConfiguration,
			TypeName: "aws_mskconnect_worker_configuration",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.KafkaConnect
}

// NewConn returns a new AWS SDK for Go v1 client for this service package's AWS API.
func (p *servicePackage) NewConn(ctx context.Context, config map[string]any) (*kafkaconnect_sdkv1.KafkaConnect, error) {
	sess := config["session"].(*session_sdkv1.Session)

	return kafkaconnect_sdkv1.New(sess.Copy(&aws_sdkv1.Config{Endpoint: aws_sdkv1.String(config["endpoint"].(string))})), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
