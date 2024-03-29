// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package location

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"terraform-provider-awsgps/internal/conns"
	"terraform-provider-awsgps/internal/errs/sdkdiag"
	tftags "terraform-provider-awsgps/internal/tags"
)

// @SDKDataSource("aws_location_route_calculator")
func DataSourceRouteCalculator() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceRouteCalculatorRead,

		Schema: map[string]*schema.Schema{
			"calculator_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"calculator_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
		},
	}
}

func dataSourceRouteCalculatorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).LocationConn(ctx)

	out, err := findRouteCalculatorByName(ctx, conn, d.Get("calculator_name").(string))
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Location Service Route Calculator (%s): %s", d.Get("calculator_name").(string), err)
	}

	if out == nil {
		return sdkdiag.AppendErrorf(diags, "reading Location Service Route Calculator (%s): empty response", d.Get("calculator_name").(string))
	}

	d.SetId(aws.StringValue(out.CalculatorName))
	d.Set("calculator_arn", out.CalculatorArn)
	d.Set("calculator_name", out.CalculatorName)
	d.Set("create_time", aws.TimeValue(out.CreateTime).Format(time.RFC3339))
	d.Set("data_source", out.DataSource)
	d.Set("description", out.Description)
	d.Set("update_time", aws.TimeValue(out.UpdateTime).Format(time.RFC3339))

	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	if err := d.Set("tags", KeyValueTags(ctx, out.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return sdkdiag.AppendErrorf(diags, "listing tags for Location Service Route Calculator (%s): %s", d.Id(), err)
	}

	return diags
}
