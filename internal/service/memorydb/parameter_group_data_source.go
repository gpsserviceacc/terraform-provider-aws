// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package memorydb

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-awsgps/internal/conns"
	"terraform-provider-awsgps/internal/errs/sdkdiag"
	tftags "terraform-provider-awsgps/internal/tags"
	"terraform-provider-awsgps/internal/tfresource"
)

// @SDKDataSource("aws_memorydb_parameter_group")
func DataSourceParameterGroup() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceParameterGroupRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"family": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parameter": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: ParameterHash,
			},
			"tags": tftags.TagsSchemaComputed(),
		},
	}
}

func dataSourceParameterGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).MemoryDBConn(ctx)
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	name := d.Get("name").(string)

	group, err := FindParameterGroupByName(ctx, conn, name)

	if err != nil {
		return sdkdiag.AppendFromErr(diags, tfresource.SingularDataSourceFindError("MemoryDB Parameter Group", err))
	}

	d.SetId(aws.StringValue(group.Name))

	d.Set("arn", group.ARN)
	d.Set("description", group.Description)
	d.Set("family", group.Family)
	d.Set("name", group.Name)

	userDefinedParameters := createUserDefinedParameterMap(d)

	parameters, err := listParameterGroupParameters(ctx, conn, d.Get("family").(string), d.Id(), userDefinedParameters)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "listing parameters for MemoryDB Parameter Group (%s): %s", d.Id(), err)
	}

	if err := d.Set("parameter", flattenParameters(parameters)); err != nil {
		return sdkdiag.AppendErrorf(diags, "failed to set parameter: %s", err)
	}

	tags, err := listTags(ctx, conn, d.Get("arn").(string))

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "listing tags for MemoryDB Parameter Group (%s): %s", d.Id(), err)
	}

	if err := d.Set("tags", tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting tags: %s", err)
	}

	return diags
}
