// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package codecatalyst

import (
	"context"

	aws_sdkv2 "github.com/aws/aws-sdk-go-v2/aws"
	codecatalyst_sdkv2 "github.com/aws/aws-sdk-go-v2/service/codecatalyst"
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
			Factory:  DataSourceDevEnvironment,
			TypeName: "aws_codecatalyst_dev_environment",
			Name:     "Dev Environment",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceDevEnvironment,
			TypeName: "aws_codecatalyst_dev_environment",
			Name:     "DevEnvironment",
		},
		{
			Factory:  ResourceProject,
			TypeName: "aws_codecatalyst_project",
			Name:     "Project",
		},
		{
			Factory:  ResourceSourceRepository,
			TypeName: "aws_codecatalyst_source_repository",
			Name:     "Source Repository",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.CodeCatalyst
}

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*codecatalyst_sdkv2.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws_sdkv2.Config))

	return codecatalyst_sdkv2.NewFromConfig(cfg, func(o *codecatalyst_sdkv2.Options) {
		if endpoint := config["endpoint"].(string); endpoint != "" {
			o.BaseEndpoint = aws_sdkv2.String(endpoint)
		}
	}), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
