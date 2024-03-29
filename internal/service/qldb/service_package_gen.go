// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package qldb

import (
	"context"

	aws_sdkv2 "github.com/aws/aws-sdk-go-v2/aws"
	qldb_sdkv2 "github.com/aws/aws-sdk-go-v2/service/qldb"
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
			Factory:  dataSourceLedger,
			TypeName: "aws_qldb_ledger",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  resourceLedger,
			TypeName: "aws_qldb_ledger",
			Name:     "Ledger",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  resourceStream,
			TypeName: "aws_qldb_stream",
			Name:     "Stream",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.QLDB
}

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*qldb_sdkv2.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws_sdkv2.Config))

	return qldb_sdkv2.NewFromConfig(cfg, func(o *qldb_sdkv2.Options) {
		if endpoint := config["endpoint"].(string); endpoint != "" {
			o.BaseEndpoint = aws_sdkv2.String(endpoint)
		}
	}), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
