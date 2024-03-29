// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package firehose

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-awsgps/internal/conns"
	"terraform-provider-awsgps/internal/errs/sdkdiag"
)

// @SDKDataSource("aws_kinesis_firehose_delivery_stream")
func dataSourceDeliveryStream() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceDeliveryStreamRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceDeliveryStreamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).FirehoseClient(ctx)

	sn := d.Get("name").(string)
	output, err := findDeliveryStreamByName(ctx, conn, sn)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Kinesis Firehose Delivery Stream (%s): %s", sn, err)
	}

	d.SetId(aws.ToString(output.DeliveryStreamARN))
	d.Set("arn", output.DeliveryStreamARN)
	d.Set("name", output.DeliveryStreamName)

	return diags
}
