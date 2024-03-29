// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package athena

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-awsgps/internal/conns"
	"terraform-provider-awsgps/internal/errs/sdkdiag"
	tfslices "terraform-provider-awsgps/internal/slices"
	"terraform-provider-awsgps/internal/tfresource"
)

// @SDKDataSource("aws_athena_named_query")
func dataSourceNamedQuery() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceNamedQueryRead,

		Schema: map[string]*schema.Schema{
			"database": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"querystring": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workgroup": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "primary",
			},
		},
	}
}

func dataSourceNamedQueryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).AthenaClient(ctx)

	input := &athena.ListNamedQueriesInput{
		WorkGroup: aws.String(d.Get("workgroup").(string)),
	}
	var queryIDs []string
	pages := athena.NewListNamedQueriesPaginator(conn, input)
	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "listing Athena Named Queries: %s", err)
		}

		queryIDs = append(queryIDs, page.NamedQueryIds...)
	}

	name := d.Get("name").(string)
	query, err := findNamedQueryByName(ctx, conn, queryIDs, name)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Athena Named Query (%s): %s", name, err)
	}

	d.SetId(aws.ToString(query.NamedQueryId))
	d.Set("database", query.Database)
	d.Set("description", query.Description)
	d.Set("name", query.Name)
	d.Set("querystring", query.QueryString)
	d.Set("workgroup", query.WorkGroup)

	return diags
}

func findNamedQueryByName(ctx context.Context, conn *athena.Client, queryIDs []string, name string) (*types.NamedQuery, error) {
	input := &athena.BatchGetNamedQueryInput{
		NamedQueryIds: queryIDs,
	}

	output, err := conn.BatchGetNamedQuery(ctx, input)

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	queries := tfslices.Filter(output.NamedQueries, func(v types.NamedQuery) bool {
		return aws.ToString(v.Name) == name
	})

	return tfresource.AssertSingleValueResult(queries)
}
