// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package licensemanager

import (
	"context"

	aws_sdkv1 "github.com/aws/aws-sdk-go/aws"
	session_sdkv1 "github.com/aws/aws-sdk-go/aws/session"
	licensemanager_sdkv1 "github.com/aws/aws-sdk-go/service/licensemanager"
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
			Factory:  DataSourceDistributedGrants,
			TypeName: "aws_licensemanager_grants",
		},
		{
			Factory:  DataSourceReceivedLicense,
			TypeName: "aws_licensemanager_received_license",
		},
		{
			Factory:  DataSourceReceivedLicenses,
			TypeName: "aws_licensemanager_received_licenses",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceAssociation,
			TypeName: "aws_licensemanager_association",
		},
		{
			Factory:  ResourceGrant,
			TypeName: "aws_licensemanager_grant",
		},
		{
			Factory:  ResourceGrantAccepter,
			TypeName: "aws_licensemanager_grant_accepter",
		},
		{
			Factory:  ResourceLicenseConfiguration,
			TypeName: "aws_licensemanager_license_configuration",
			Name:     "License Configuration",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "id",
			},
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.LicenseManager
}

// NewConn returns a new AWS SDK for Go v1 client for this service package's AWS API.
func (p *servicePackage) NewConn(ctx context.Context, config map[string]any) (*licensemanager_sdkv1.LicenseManager, error) {
	sess := config["session"].(*session_sdkv1.Session)

	return licensemanager_sdkv1.New(sess.Copy(&aws_sdkv1.Config{Endpoint: aws_sdkv1.String(config["endpoint"].(string))})), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
