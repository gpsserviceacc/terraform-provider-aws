// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package secretsmanager

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-awsgps/internal/conns"
	"terraform-provider-awsgps/internal/errs/sdkdiag"
	"terraform-provider-awsgps/internal/generate/namevaluesfiltersv2"
	tfslices "terraform-provider-awsgps/internal/slices"
)

// @SDKDataSource("aws_secretsmanager_secrets", name="Secrets")
func dataSourceSecrets() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceSecretsRead,
		Schema: map[string]*schema.Schema{
			"arns": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"filter": namevaluesfiltersv2.Schema(),
			"names": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceSecretsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SecretsManagerClient(ctx)

	input := &secretsmanager.ListSecretsInput{}

	if v, ok := d.GetOk("filter"); ok {
		input.Filters = namevaluesfiltersv2.New(v.(*schema.Set)).SecretsmanagerFilters()
	}

	var results []types.SecretListEntry

	paginator := secretsmanager.NewListSecretsPaginator(conn, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "listing Secrets Manager Secrets: %s", err)
		}

		if page != nil {
			results = append(results, page.SecretList...)
		}
	}

	d.SetId(meta.(*conns.AWSClient).Region)
	d.Set("arns", tfslices.ApplyToAll(results, func(v types.SecretListEntry) string { return aws.ToString(v.ARN) }))
	d.Set("names", tfslices.ApplyToAll(results, func(v types.SecretListEntry) string { return aws.ToString(v.Name) }))

	return diags
}
