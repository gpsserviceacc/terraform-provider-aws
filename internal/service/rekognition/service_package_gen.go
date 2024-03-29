// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package rekognition

import (
	"context"

	aws_sdkv2 "github.com/aws/aws-sdk-go-v2/aws"
	rekognition_sdkv2 "github.com/aws/aws-sdk-go-v2/service/rekognition"
	"terraform-provider-awsgps/internal/conns"
	"terraform-provider-awsgps/internal/types"
	"terraform-provider-awsgps/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{
		{
			Factory: newResourceCollection,
			Name:    "Collection",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory: newResourceProject,
			Name:    "Project",
		},
	}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{}
}

func (p *servicePackage) ServicePackageName() string {
	return names.Rekognition
}

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*rekognition_sdkv2.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws_sdkv2.Config))

	return rekognition_sdkv2.NewFromConfig(cfg, func(o *rekognition_sdkv2.Options) {
		if endpoint := config["endpoint"].(string); endpoint != "" {
			o.BaseEndpoint = aws_sdkv2.String(endpoint)
		}
	}), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
