// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package codeartifact

import (
	"context"

	aws_sdkv2 "github.com/aws/aws-sdk-go-v2/aws"
	codeartifact_sdkv2 "github.com/aws/aws-sdk-go-v2/service/codeartifact"
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
			Factory:  dataSourceAuthorizationToken,
			TypeName: "aws_codeartifact_authorization_token",
			Name:     "Authoiration Token",
		},
		{
			Factory:  dataSourceRepositoryEndpoint,
			TypeName: "aws_codeartifact_repository_endpoint",
			Name:     "Repository Endpoint",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  resourceDomain,
			TypeName: "aws_codeartifact_domain",
			Name:     "Domain",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  resourceDomainPermissionsPolicy,
			TypeName: "aws_codeartifact_domain_permissions_policy",
			Name:     "Domain Permissions Policy",
		},
		{
			Factory:  resourceRepository,
			TypeName: "aws_codeartifact_repository",
			Name:     "Repository",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  resourceRepositoryPermissionsPolicy,
			TypeName: "aws_codeartifact_repository_permissions_policy",
			Name:     "Repository Permissions Policy",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.CodeArtifact
}

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*codeartifact_sdkv2.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws_sdkv2.Config))

	return codeartifact_sdkv2.NewFromConfig(cfg, func(o *codeartifact_sdkv2.Options) {
		if endpoint := config["endpoint"].(string); endpoint != "" {
			o.BaseEndpoint = aws_sdkv2.String(endpoint)
		}
	}), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
