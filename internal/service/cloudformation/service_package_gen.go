// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package cloudformation

import (
	"context"

	aws_sdkv1 "github.com/aws/aws-sdk-go/aws"
	session_sdkv1 "github.com/aws/aws-sdk-go/aws/session"
	cloudformation_sdkv1 "github.com/aws/aws-sdk-go/service/cloudformation"
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
			Factory:  DataSourceExport,
			TypeName: "aws_cloudformation_export",
		},
		{
			Factory:  DataSourceStack,
			TypeName: "aws_cloudformation_stack",
		},
		{
			Factory:  DataSourceType,
			TypeName: "aws_cloudformation_type",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceStack,
			TypeName: "aws_cloudformation_stack",
			Name:     "Stack",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  ResourceStackSet,
			TypeName: "aws_cloudformation_stack_set",
			Name:     "Stack Set",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  ResourceStackSetInstance,
			TypeName: "aws_cloudformation_stack_set_instance",
		},
		{
			Factory:  ResourceType,
			TypeName: "aws_cloudformation_type",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.CloudFormation
}

// NewConn returns a new AWS SDK for Go v1 client for this service package's AWS API.
func (p *servicePackage) NewConn(ctx context.Context, config map[string]any) (*cloudformation_sdkv1.CloudFormation, error) {
	sess := config["session"].(*session_sdkv1.Session)

	return cloudformation_sdkv1.New(sess.Copy(&aws_sdkv1.Config{Endpoint: aws_sdkv1.String(config["endpoint"].(string))})), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
